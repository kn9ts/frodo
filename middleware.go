package frodo

import (
	"net/http"
	"strings"
)

// Request will help facilitate the passing of multiple handlers
type Request struct {
	handlers       []Handle
	total          int
	nextPosition   int
	files          []*UploadedFile
	ResponseWriter *ResponseWriter
	*http.Request
	Params
}

// Middleware declares the minimum implementation necessary for a handlers
// to be used as Frodo's middleware route Handlers
type Middleware interface {
	Next(w *ResponseWriter)
}

func (r *Request) runHandleChain(w *ResponseWriter) {
	r.nextPosition++
	r.ResponseWriter = w
	r.handlers[0](w, r)
}

// Next will be used to call the next handler in line/queue
func (r *Request) Next() {
	// 1st check if the next handler position accounts for the number
	// of handlers existing in the handlers array
	if r.nextPosition < r.total {
		// get the next handler
		nextHandler := r.handlers[r.nextPosition]
		// move the cursor
		r.nextPosition++

		// 1st check if a write has happened
		// meaning a response has been issued out to the client
		// if not run the next handler in line
		if r.ResponseWriter.ResponseSent() == false {
			nextHandler(r.ResponseWriter, r)
		}
	}
	return
}

// Input gets ALL key/values sent via POST from all methods.
// Keep in mind `r.Form == type url.Values map[string][]string`
func (r *Request) Input(name string) interface{} {
	if r.Form == nil {
		r.ParseForm()
	}

	if value, ok := r.Form[name]; ok {
		if len(value) == 1 {
			return value[0]
		}
		return value
	}

	return nil
}

// HasInput checks for the existence of the given
// input name in the inputs sent from a FORM
func (r *Request) HasInput(name string) bool {
	if r.Form == nil {
		r.ParseForm()
	}

	_, ok := r.Form[name]
	return ok
}

// HasFile mimics FormFile method from `http.Request`
//      func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)
func (r *Request) HasFile(name string) bool {
	_, _, err := r.FormFile(name)
	if err != nil {
		return false
	}
	return true
}

// UploadedFile gets the file requested that was uploaded
func (r *Request) UploadedFile(name string) (*UploadedFile, error) {
	file, header, err := r.FormFile(name)
	if err == nil {
		return &UploadedFile{file, header}, nil
	}
	return nil, err
}

// UploadedFiles parses all uploaded files and creates an
// array of UploadedFile type representing each uploaded file
func (r *Request) UploadedFiles(name string) []*UploadedFile {
	// Instantiate r.files
	if r.files == nil {
		r.files = make([]*UploadedFile, len(r.MultipartForm.File[name]))
		r.ParseMultipartForm(32 << 20)
	}

	for _, header := range r.MultipartForm.File[name] {
		file, _ := header.Open()
		r.files = append(r.files, &UploadedFile{file, header})
	}

	return r.files
}

// MoveAll is a neat trick to upload all the files that
// have been parsed. Awesome for bulk uploading, and storage.
func (r *Request) MoveAll(args ...interface{}) (bool, int) {
	if r.files == nil {
		return false, 0
	}

	count := 0
	for _, file := range r.files {
		moved := file.Move(args...)
		if moved {
			count++
		}
	}

	if count == len(r.files) {
		return true, count
	}
	return false, count
}

// ClientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
func (r *Request) ClientIP() string {
	if true {
		clientIP := strings.TrimSpace(r.Request.Header.Get("X-Real-Ip"))
		if len(clientIP) > 0 {
			return clientIP
		}
		clientIP = r.Request.Header.Get("X-Forwarded-For")
		if index := strings.IndexByte(clientIP, ','); index >= 0 {
			clientIP = clientIP[0:index]
		}
		clientIP = strings.TrimSpace(clientIP)
		if len(clientIP) > 0 {
			return clientIP
		}
	}
	return strings.TrimSpace(r.Request.RemoteAddr)
}

// IsAjax checks if the Request was made via AJAX,
// the XMLHttpRequest will usually be sent with a X-Requested-With HTTP header.
func (r *Request) IsAjax() bool {
	if r.Request.Header.Get("X-Request-With") != "" {
		return true
	}
	return false
}

// IsXhr gives user a choice in whichever way he/she feels okay checking for AJAX Request
// It actually calls r.IsAjax()
func (r *Request) IsXhr() bool {
	return r.IsAjax()
}

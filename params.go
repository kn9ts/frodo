package Frodo

import "net/http"

// Request type will carry all the http.Request values, and params in curly {} brackets that are
// translated from url param values to ready to be used values
// also decided to handle file uploads
type Request struct {
	*http.Request
	params, form map[string]string

	// To cater for multiple file uploads
	files []*UploadFile
}

// GetParam returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (r *Request) GetParam(name string) string {
	value, ok := r.params[name]
	if ok {
		return value
	}
	return ""
}

// Param is equalin using get
func (r *Request) Param(name string) string {
	return r.GetParam(name)
}

// SetParam adds a key/value pair to the Request params
func (r *Request) SetParam(name, value string) bool {
	// 1st check if it has been initialised
	if r.params != nil { // If not initialise
		r.params = make(map[string]string)
	}

	// allow overwriting
	r.params[name] = value
	return true
}

// Input gets ALL posted key/values from all Methods
// Keep in mind r.Form == type url.Values map[string][]string
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

// HasInput checks for the existence of the given input name in the inputs
func (r *Request) HasInput(name string) bool {
	if r.Form == nil {
		r.ParseForm()
	}

	_, ok := r.Form[name]
	return ok
}

// HasFile mimics FormFile method for Request
// ------ func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error)
func (r *Request) HasFile(name string) bool {
	_, _, err := r.FormFile(name)
	if err != nil {
		return false
	}
	return true
}

// InputFile gets the file requested that was uploaded
func (r *Request) InputFile(name string) (*UploadFile, error) {
	file, header, err := r.FormFile(name)
	if err == nil {
		return &UploadFile{file, header}, nil
	}
	return nil, err
}

// InputFiles parses all uploaded files and creates an array of UploadFile
// representing each uploaded file
func (r *Request) InputFiles(name string) []*UploadFile {
	// Instantiate r.files
	if r.files == nil {
		r.files = make([]*UploadFile, len(r.MultipartForm.File[name]))
		r.ParseMultipartForm(32 << 20)
	}

	for _, header := range r.MultipartForm.File[name] {
		file, _ := header.Open()
		r.files = append(r.files, &UploadFile{file, header})
	}

	return r.files
}

// MoveAll is a neat trick for uploads all the files that have been parsed
// Awesome for bulk uploading, and storage
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

// IsAjax checks if the Request was made via AJAX
// the XMLHttpRequest will usually be sent with a X-Requested-With HTTP header.
func (r *Request) IsAjax() bool {
	if r.Header.Get("X-Request-With") != "" {
		return true
	}
	return false
}

// IsXhr gives user a choice in whichever way he/she feels okay checking for AJAX Request
// It actually calls r.IsAjax()
func (r *Request) IsXhr() bool {
	return r.IsAjax()
}

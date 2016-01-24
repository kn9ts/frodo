package frodo

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"time"

	"log"
)

const errorMessage string = "[ERROR] Headers were already written."

// ResponseWriter is used to hijack/embed http.ResponseWriter
// thus making it satisfy the ResponseWriter interface, we then add a written boolean property
// to trace when a write made, with a couple of other helpful properties
type ResponseWriter struct {
	http.ResponseWriter
	headerWritten bool
	written       bool
	timeStart     time.Time
	timeEnd       time.Time
	duration      float64
	statusCode    int
	size          int64
	method        string
	route         string
}

// Write writes data back the client/creates the body
func (w *ResponseWriter) Write(bytes []byte) (int, error) {
	if !w.HeaderWritten() {
		w.WriteHeader(http.StatusOK)
	}

	if w.ResponseSent() {
		log.Println(errorMessage)
		return 1, errors.New(errorMessage)
	}

	sent, err := w.ResponseWriter.Write(bytes)
	if err != nil {
		return sent, err
	}
	w.size += int64(sent)
	w.timeEnd = time.Now()
	w.duration = time.Since(w.timeEnd).Seconds()
	return sent, nil
}

// WriteHeader writes the Headers out
func (w *ResponseWriter) WriteHeader(code int) {
	if w.HeaderWritten() {
		log.Println(errorMessage)
		return
	}
	w.ResponseWriter.WriteHeader(code)
	w.headerWritten = true
	w.statusCode = code
}

// ResponseSent checks if a write has been made
// starts with a header being sent out
func (w *ResponseWriter) ResponseSent() bool {
	return w.written
}

// HeaderWritten checks if a write has been made
// starts with a header being sent out
func (w *ResponseWriter) HeaderWritten() bool {
	return w.headerWritten
}

// Flush wraps response writer's Flush function.
func (w *ResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

// Hijack wraps response writer's Hijack function.
func (w *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// CloseNotify wraps response writer's CloseNotify function.
func (w *ResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Size returns the size of the response
// about to be sent out
func (w *ResponseWriter) Size() int64 {
	return w.size
}

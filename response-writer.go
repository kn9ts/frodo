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

// AppResponseWriter is used to hijack/embed http.ResponseWriter
// thus making it satisfy the ResponseWriter interface, we then add a written boolean property
// to trace when a write made, with a couple of other helpful properties
type AppResponseWriter struct {
	http.ResponseWriter
	written   bool
	timeStart time.Time
	timeEnd   time.Time
	duration  time.Duration
	status    int
	size      int64
	method    string
	route     string
}

// Write writes data back the client/creates the body
func (w *AppResponseWriter) Write(bytes []byte) (int, error) {
	if !w.WriteHappened() {
		w.WriteHeader(http.StatusOK)
	} else {
		log.Println(errorMessage)
		return 1, errors.New(errorMessage)
	}

	done, err := w.ResponseWriter.Write(bytes)
	if err != nil {
		return done, err
	}
	w.size += int64(done)
	w.timeEnd = time.Now()
	w.duration = time.Since(w.timeEnd)
	return done, nil
}

// WriteHeader writes the Headers out
func (w *AppResponseWriter) WriteHeader(code int) {
	if w.WriteHappened() {
		log.Println(errorMessage)
		return
	}
	w.ResponseWriter.WriteHeader(code)
	w.status = code
	w.written = true
}

// WriteHappened checks if a write has been made
// starts with a header being sent out
func (w *AppResponseWriter) WriteHappened() bool {
	return w.written
}

// Flush wraps response writer's Flush function.
func (w *AppResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

// Hijack wraps response writer's Hijack function.
func (w *AppResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// CloseNotify wraps response writer's CloseNotify function.
func (w *AppResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Status gets the status code of the response
func (w *AppResponseWriter) Status() int {
	return w.status
}

// Size returns the size of the response
// about to be sent out
func (w *AppResponseWriter) Size() int64 {
	return w.size
}

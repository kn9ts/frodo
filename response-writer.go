package frodo

import (
	"bufio"
	"net"
	"net/http"
	"time"

	"log"
)

// WrappedResponseWriter is used to hijack/embed http.ResponseWriter
// thus making it satisfy the ResponseWriter interface, we then add a written boolean property
// to trace when a write made, with a couple of other helpful properties
type WrappedResponseWriter struct {
	http.ResponseWriter
	written   bool
	timeStart time.Time
	timeEnd   time.Time
	duration  float64
	status    int
	size      int64
}

// Header returns the response header
// use it if one desires to change or add header info to the reponse
// before it is sent out/back
func (r *WrappedResponseWriter) Header() http.Header {
	return r.ResponseWriter.Header()
}

// Write writes data back the client/creates the body
func (r *WrappedResponseWriter) Write(bytes []byte) (int, error) {
	r.WriteHeader(http.StatusOK)
	n, err := r.ResponseWriter.Write(bytes)
	if err != nil {
		return n, err
	}
	r.size += int64(n)
	return n, nil
}

// WriteHeader writes the Headers out
func (r *WrappedResponseWriter) WriteHeader(code int) {
	if r.written {
		log.Println("[WARNING] Headers were already written. Can not write out another one.")
		return
	}
	r.status = code
	r.WriteHeader(code)
	r.written = true
}

// WriteHappened checks if a write has been made
// starts with a header being sent out
func (r *WrappedResponseWriter) WriteHappened() bool {
	return r.written
}

// Flush wraps response writer's Flush function.
func (r *WrappedResponseWriter) Flush() {
	r.ResponseWriter.(http.Flusher).Flush()
}

// Hijack wraps response writer's Hijack function.
func (r *WrappedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.ResponseWriter.(http.Hijacker).Hijack()
}

// CloseNotify wraps response writer's CloseNotify function.
func (r *WrappedResponseWriter) CloseNotify() <-chan bool {
	return r.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Status gets the status code of the response
func (r *WrappedResponseWriter) Status() int {
	return r.status
}

// Size returns the size of the response
// about to be sent out
func (r *WrappedResponseWriter) Size() int64 {
	return r.size
}

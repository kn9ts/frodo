package Frodo

import (
	"fmt"
	"net/http"
)

// MiddlewareResponseWriter is used to hijack/embed http.ResponseWriter
// thus making it satisfy the ResponseWriter interface, we then add a written boolean property
// to trace when a write made and exit
type MiddlewareResponseWriter struct {
	http.ResponseWriter
	written bool
}

// Write writes data back the client/creates the body
func (w *MiddlewareResponseWriter) Write(bytes []byte) (int, error) {
	w.written = true
	fmt.Printf("\nAn application response was written back: %v\n", w.written)
	return w.ResponseWriter.Write(bytes)
}

// WriteHeader is in charge of building the Header file and writing it back to the client
func (w *MiddlewareResponseWriter) WriteHeader(code int) {
	w.written = true
	fmt.Printf("\nHeader was wriiten back: %v\n", w.written)
	w.ResponseWriter.WriteHeader(code)
}

/*
   FUTURE:

   If a write happens, then we can track it here and we can place our
   AFTER middleware just before a write happens
*/

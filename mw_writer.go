package Frodo

import (
	"fmt"
	"net/http"
)

type MiddlewareResponseWriter struct {
	http.ResponseWriter
	written bool
}

func (w *MiddlewareResponseWriter) Write(bytes []byte) (int, error) {
	w.written = true
	fmt.Printf("\nA write was made by Middleware: %v\n", w.written)
	return w.ResponseWriter.Write(bytes)
}

func (w *MiddlewareResponseWriter) WriteHeader(code int) {
	w.written = true
	fmt.Printf("\nA write was made by Middleware: %v\n", w.written)
	w.ResponseWriter.WriteHeader(code)
}

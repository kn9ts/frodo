package frodo

import "net/http"

// Request will help facilitate the passing of multiple handlers
type Request struct {
	handlers       []Handle
	total          int
	nextPosition   int
	ResponseWriter *AppResponseWriter
	*http.Request
	Params
}

// Middleware declares the minimum implementation necessary for a handlers
// to be used as Frodo's middleware route Handlers
type Middleware interface {
	Next(w *AppResponseWriter)
}

func (r *Request) runHandleChain(w *AppResponseWriter) {
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
		if r.ResponseWriter.WriteHappened() == false {
			nextHandler(r.ResponseWriter, r)
		}
	}
	return
}

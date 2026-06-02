package frodo

import "net/http"

// Middleware declares the minimum implementation
// necessary for a handlers to be used as Frodo's middleware route Handlers
type Middleware interface {
	Next(...interface{})
}

// RequestMiddleware is a collection of Middlwares
// this struct will have the method next to invoke the middleware in the chain
// one after the other
type RequestMiddleware struct {
	handlers       []Middleware // []Handle
	ResponseWriter *ResponseWriter
	Request        *Request
	total          int
	nextPosition   int
}

func (m *RequestMiddleware) chainReaction() {
	m.nextPosition++
	m.typeCastAndCall(m.handlers[m.nextPosition-1])
}

// Next will be used to call the next handler in line/queue
func (m *RequestMiddleware) Next(args ...interface{}) {
	// 1st check if the next handler position accounts for the number
	// of handlers existing in the handlers array
	if m.nextPosition < m.total {
		// get the next handler
		nextHandle := m.handlers[m.nextPosition]
		// move the cursor
		m.nextPosition++
		// 1st check if a write has happened
		// have headers been sent out?
		// meaning a response has been issued out to the client
		// if not call the next handler in line
		if !m.ResponseWriter.HeaderWritten() {
			m.typeCastAndCall(nextHandle)
		}
	}
}

// typeCastAndCall converts the middleware to it's rightful type then calls it
func (m *RequestMiddleware) typeCastAndCall(run Middleware) {
	// 1st check if the route handler is HandleFunc
	if handle, hasTypeCasted := run.(Handler); hasTypeCasted {
		handle(m.ResponseWriter, m.Request)
	} else {
		// if not, then is it an implementation of ControllerInterface
		if ctrl, ok := run.(CRUDController); ok {
			ctrl.Index(m.ResponseWriter, m.Request)
		} else {
			m.ResponseWriter.WriteHeader(404)
			http.Error(m.ResponseWriter, "No Frodo.Handle or Frodo.Controller exists to handle the route.", 404)
		}
	}
}

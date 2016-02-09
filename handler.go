package frodo

import (
	"net/http"
)

// Handler is a function that can be registered to a route to handle HTTP requests.
type Handler func(http.ResponseWriter, *Request)

// Next enables Handler types to be treated as Middleware too
func (h Handler) Next(r ...interface{}) {
}

// ControllerHandle is used to incubate the Controller Methods and it's Attributes
// since the Attributes are lost in the previous ways of attaching them to routes
type ControllerHandle struct {
	Handler    CRUDController
	Attributes Attributes
}

// Next enables ControllerHandle types to be treated as Middleware too
func (h ControllerHandle) Next(r ...interface{}) {
}

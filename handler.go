package frodo

import (
	"net/http"
)

// Handler is a function that can be registered to a route to handle HTTP requests.
type Handler func(http.ResponseWriter, *Request)

// Next enables Handler types to be treated as Middleware too
func (h Handler) Next(r ...interface{}) {
}

package frodo

import (
	"net/http"
)

// Request wraps *http.Request and httprouter.Params into one
type Request struct {
	Request *http.Request
	Params  Params
}

// Context finally wraps http.ResponseWriter, *http.Request and httprouter.Params into one
type Context struct {
	Response http.ResponseWriter
	Request
}

// ContextFunc would be an optional way of defining a route Handler
type ContextFunc func(c *Context)

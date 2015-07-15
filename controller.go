package Frodo

import (
	"net/http"
)

// Controller Defines the basic structure of a REST Application Controller
// the devs, controller should embedd this
type Controller struct {
	Method, Layout string
}

func (c *Controller) _Construct(fn func(*Controller)) {
	fn(c)
}

// Default is the fallback Method when no method has been provided in router by the Controller
func (c *Controller) Default(w http.ResponseWriter, r *Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Get adds a request function to handle GET request.
func (c *Controller) Get(w http.ResponseWriter, r *Request) {
	// http.Error(w, "Method Not Allowed", 405)
}

// Post adds a request function to handle POST request.
func (c *Controller) Post(w http.ResponseWriter, r *Request) {
	// http.Error(w, "Method Not Allowed", 405)
}

// Put adds a request function to handle PUT request.
func (c *Controller) Put(w http.ResponseWriter, r *Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Patch adds a request function to handle PATCH request.
func (c *Controller) Patch(w http.ResponseWriter, r *Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Delete adds a request function to handle DELETE request.
func (c *Controller) Delete(w http.ResponseWriter, r *Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Head adds a request function to handle HEAD request.
func (c *Controller) Head(w http.ResponseWriter, r *Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Options adds a request function to handle OPTIONS request.
func (c *Controller) Options(w http.ResponseWriter, r *Request) {
	http.Error(w, "Method Not Allowed", 405)
}

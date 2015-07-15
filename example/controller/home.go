package controller

import (
	"github.com/kn9ts/frodo"
	"net/http"
)

// Home is plays an example of a controller
type Home struct {
	Frodo.Controller
}

// Index is the default route handler for "/" route
func (h *Home) Index(w http.ResponseWriter, r *Frodo.Request) {
	w.Write([]byte("Hello world"))
}

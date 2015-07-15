package controller

import (
	"github.com/kn9ts/frodo"
	"net/http"
)

// Home is plays an example of a controller
type Home struct {
	Frodo.Controller
}

// Get is overriding and defining your own get method
func (h *Home) Get(w http.ResponseWriter, r *Frodo.Request) {
	w.Write([]byte("Hello world"))
}

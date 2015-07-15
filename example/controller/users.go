package controller

import (
	"github.com/kn9ts/frodo"
	"net/http"
)

// Users Controller instance
type Users struct {
	Frodo.Controller
}

// Get is overriding and defining your own get method
func (u *Users) Get(w http.ResponseWriter, r *Frodo.Request) {
	w.Write([]byte("Hello a list of users will be here"))
}

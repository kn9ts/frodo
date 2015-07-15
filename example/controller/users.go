package controller

import (
	"github.com/kn9ts/frodo"
	"net/http"
)

// Users Controller instance
type Users struct {
	Frodo.Controller
}

// Create is overriding and defining your Create method
func (u *Users) Create(w http.ResponseWriter, r *Frodo.Request) {
	w.Write([]byte("Hello a list of users will be here"))
}

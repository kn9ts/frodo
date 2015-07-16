package controller

import (
	"github.com/kn9ts/frodo"
	"net/http"
)

// Users Controller instance
type Users struct {
	Frodo.Controller
}

// Show is overriding and defining user's Show method
func (u *Users) Show(w http.ResponseWriter, r *Frodo.Request) {
	w.Write([]byte("Hello a list of users will be here"))
}

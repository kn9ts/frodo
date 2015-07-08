package main

import "../../frodo"
import "net/http"

func main() {
	App := Frodo.New.Application()
	// Now create your routes
	App.Get("/", func(w http.ResponseWriter, r *http.Request, _ Frodo.Params) {
		w.Write([]byte("Welcome, to the ROOT route")) // send data to client side
	})

	App.Get("/page/{id}", func(w http.ResponseWriter, r *http.Request, params Frodo.Params) {
		w.Write([]byte("Hello page here, the ID passed is " + params.Get("id"))) // send data to client side
	})

	App.Post("/{some_id}/{images}", func(w http.ResponseWriter, r *http.Request, input Frodo.Params) {
		w.Write([]byte("Hello, to get here. You required this ID: " + input.Get("some_id"))) // send data to client side
	})

	App.Post("/{name}", func(w http.ResponseWriter, r *http.Request, input Frodo.Params) {
		w.Write([]byte("Hello, " + input.Get("name") + "! This is your profile page.")) // send data to client side
	})
	App.Run(4500)
}

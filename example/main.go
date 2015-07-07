package main

import "github.com/kn9ts/frodo"
import "net/http"

func main() {
	App := frodo.New.Application()
	// Now create your routes
	App.Get("/", func(w http.ResponseWriter, r *http.Request, _ frodo.Params) {
		w.Write([]byte("Welcome, to the ROOT route")) // send data to client side
	})

	App.Get("/page/{id}", func(w http.ResponseWriter, r *http.Request, _ frodo.Params) {
		w.Write([]byte("Hello page with ID!")) // send data to client side
	})

	App.Post("/{name}", func(w http.ResponseWriter, r *http.Request, _ frodo.Params) {
		w.Write([]byte("Hello nested page that accepts names as param!")) // send data to client side
	})

	App.Post("/page", func(w http.ResponseWriter, r *http.Request, _ frodo.Params) {
		w.Write([]byte("Hello nested page called page!")) // send data to client side
	})
	App.Run()
}

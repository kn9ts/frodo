package main

import "github.com/kn9ts/frodo"
import "net/http"

func main() {
	App := frodo.New.Application()

	// Now create your routes
	App.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome, visit sub pages now.")) // send data to client side
	})

	App.Get("/page", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello nested page!")) // send data to client side
	})

	App.Serve()
}

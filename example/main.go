package main

import "github.com/kn9ts/frodo"
import "net/http"
import "fmt"

func main() {
	App := frodo.New.Application()

	// Now create your routes
	App.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!") // send data to client side
	})

	App.Get("/home", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Eugene!") // send data to client side
	})

	App.Serve()
}

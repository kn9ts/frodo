package main

import (
	"github.com/kn9ts/frodo"
	"github.com/kn9ts/frodo/example/filters"
	"gopkg.in/unrolled/render.v1"
	"net/http"
)

func main() {
	// Get an app instance by initializing the app
	App := Frodo.New.Application()
	App.AddFilters(filters.MiddleWare)
	Reponse := render.New(render.Options{})

	// Now create your routes
	App.Get("/", func(w http.ResponseWriter, r *http.Request, _ Frodo.Params) {
		w.Write([]byte("Welcome, to the ROOT route")) // send data to client side
	})

	App.Get("/page/{id}", func(w http.ResponseWriter, r *http.Request, params Frodo.Params) {
		w.Write([]byte("Hello page here, the ID passed is " + params.Get("id"))) // send data to client side
	})

	App.Post("/{some_id}/{images}", func(w http.ResponseWriter, r *http.Request, param Frodo.Params) {
		w.Write([]byte("Hello, to get here. You required this ID: " + param.Get("some_id"))) // send data to client side
	})

	App.Post("/{name}", func(w http.ResponseWriter, r *http.Request, _ Frodo.Params) {
		// w.Write([]byte("Hello, " + param.Get("name") + "! This is your profile page.")) // send data to client side
		Reponse.JSON(w, http.StatusOK, r)
	})

	App.Match(Frodo.Methods{"GET", "POST"}, "/home", func(w http.ResponseWriter, r *http.Request, _ Frodo.Params) {
		Reponse.JSON(w, http.StatusOK, r)
	})

	App.All("/home", func(w http.ResponseWriter, r *http.Request, _ Frodo.Params) {
		Reponse.JSON(w, http.StatusOK, r)
	})

	App.Serve()
	// App.ServeOnPort(3000)
}

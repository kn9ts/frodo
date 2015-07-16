package main

import (
	"github.com/kn9ts/frodo"
	"github.com/kn9ts/frodo/example/controller"
	"github.com/kn9ts/frodo/example/filters"
	"gopkg.in/unrolled/render.v1"
	"net/http"
)

func main() {
	// Get an app instance by initializing the app
	App := Frodo.New.Application()
	Reponse := render.New(render.Options{})

	// route, handler/controller, which method for controller, name of route if you plan to give it one
	// If no Method is given, then Default method kicks or decides via the HTTP Method channeling to the Controller
	App.All("/", func(w http.ResponseWriter, r *Frodo.Request) {
		w.Write([]byte("Hello World"))
	})

	App.Get("/me", &controller.Users{})
	App.Get("/users", &controller.Users{}, Frodo.Use{Method: "Show", Name: "users", Filter: "auth"})
	App.Post("/profile", &controller.Home{}, Frodo.Use{Method: "Index", Name: "profile"})

	App.Get("/settings", func(w http.ResponseWriter, r *Frodo.Request) {
		w.Write([]byte("Hello Setting"))
	}, Frodo.Use{Method: "", Name: "settings"})

	App.Get("/page/{id}", func(w http.ResponseWriter, r *Frodo.Request) {
		w.Write([]byte("Hello page here, the ID passed is " + r.Param("id"))) // send data to client side
	})

	App.Post("/payments", func(w http.ResponseWriter, r *Frodo.Request) {
		w.Write([]byte("Hello Payments"))
	}, "payments")

	App.Match(Frodo.Methods{"GET", "POST"}, "/home", func(w http.ResponseWriter, r *Frodo.Request) {
		Reponse.JSON(w, http.StatusOK, r)
	})

	App.Post("/{name}", func(w http.ResponseWriter, r *Frodo.Request) {
		// send data to client side
		w.Write([]byte("Hello, " + r.Param("name") + "! This is your profile page."))
		// Reponse.JSON(w, http.StatusOK, r)
	})

	Frodo.Log.FilePath = "./"
	Frodo.Log.WriteToFile()

	App.AddFilters(filters.MiddleWare)
	App.Serve()
	// App.ServeOnPort(3000)
}

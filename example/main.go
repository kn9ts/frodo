package main

import (
	"github.com/kn9ts/frodo"
	"github.com/kn9ts/frodo/example/controller"
	"github.com/kn9ts/frodo/example/filters"
	// "gopkg.in/unrolled/render.v1"
	"net/http"
)

func main() {
	// Get an app instance by initializing the app
	App := Frodo.New.Application()
	// Reponse := render.New(render.Options{})

	// route, handler/controller, which method for controller, name of route if you plan to give it one
	// If no Method is given, then Default method kicks or decides via the HTTP Method channeling to the Controller
	App.Get("/", func(w http.ResponseWriter, r *Frodo.Request) {
		w.Write([]byte("Hello World"))
	})

	App.Get("/settings", func(w http.ResponseWriter, r *Frodo.Request) {
		w.Write([]byte("Hello Setting"))
	}, Frodo.Use{Method: "Index"})

	App.Get("/payments", func(w http.ResponseWriter, r *Frodo.Request) {
		w.Write([]byte("Hello Payments"))
	}, "payments")

	App.Get("/me", &controller.Users{})
	App.Get("/users", &controller.Users{}, Frodo.Use{Method: "Create", Name: "users", Before: "auth", After: "token"})
	App.Get("/profile", &controller.Home{}, Frodo.Use{Method: "Index", Name: "profile"})

	// App.Get("/page/{id}", func(w http.ResponseWriter, r *Frodo.Request) {
	// 	w.Write([]byte("Hello page here, the ID passed is " + r.Param("id"))) // send data to client side
	// })
	//
	// App.Match(Frodo.Methods{"GET", "POST"}, "/home", func(w http.ResponseWriter, r *Frodo.Request) {
	// 	Reponse.JSON(w, http.StatusOK, r)
	// })
	//
	// App.Post("/{name}", func(w http.ResponseWriter, r *Frodo.Request) {
	// 	// send data to client side
	// 	// w.Write([]byte("Hello, " + param.Get("name") + "! This is your profile page."))
	// 	Reponse.JSON(w, http.StatusOK, r)
	// })

	Frodo.Log.FilePath = "./"
	Frodo.Log.WriteToFile()

	App.AddFilters(filters.MiddleWare)
	App.Serve()
	// App.ServeOnPort(3000)
}

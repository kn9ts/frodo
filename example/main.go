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

    App.All("/", func(w http.ResponseWriter, r *Frodo.Request) {
        Reponse.JSON(w, http.StatusOK, r)
    })

    App.Match(Frodo.Methods{"GET", "POST"}, "/home", func(w http.ResponseWriter, r *Frodo.Request) {
        Reponse.JSON(w, http.StatusOK, r)
    })

    App.Get("/page/{id}", func(w http.ResponseWriter, r *Frodo.Request) {
        w.Write([]byte("Hello page here, the ID passed is " + r.Param("id"))) // send data to client side
    })

    App.Post("/{some_id}/{images}", func(w http.ResponseWriter, r *Frodo.Request) {
        // send data to client side
        w.Write([]byte("Hello, to get here. You required this ID: " + r.Param("some_id")))
    })

    App.Post("/{name}", func(w http.ResponseWriter, r *Frodo.Request) {
        // send data to client side
        // w.Write([]byte("Hello, " + param.Get("name") + "! This is your profile page."))
        Reponse.JSON(w, http.StatusOK, r)
    })

    App.Serve()
    // App.ServeOnPort(3000)
}

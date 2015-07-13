package filters

import (
    "github.com/kn9ts/frodo"
    "net/http"
)

// This is how a dev makes a Filter instance
var MiddleWare = Frodo.NewFilters()

func init() {

    // Adding Before/After filters
    MiddleWare.Before(func(w http.ResponseWriter, r *Frodo.Request) {
        if r.Method == "GET" {
            r.Method = "CHANGED_BY_BEFORE_MIDDLEWARE"
        }
        // w.Write([]byte("Middleware wrote this, so the application should exit.\n"))
    })

    MiddleWare.After(func(w http.ResponseWriter, r *Frodo.Request) {
        if r.Method == "GET" {
            r.Method = "CHANGED_BY_AFTER_MIDDLEWARE"
        }
    })

    // Adding routin filters, this applies now to "/page/{id}" route
    MiddleWare.Filter("/page/{id}", func(w http.ResponseWriter, r *Frodo.Request) {
        if r.GetParam("id") != "" {
            r.Method = "CHANGED_BY_FILTER_MIDDLEWARE"
        }
    }, false)

}

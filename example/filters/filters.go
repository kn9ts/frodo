package filters

import (
	"github.com/kn9ts/frodo"
	"net/http"
)

// This is how a dev makes a Filter instance
var MiddleWare = Frodo.NewFilters()

func init() {

	// Adding Before/After filters
	MiddleWare.Before(func(w http.ResponseWriter, r *http.Request, _ Frodo.Params) {
		if r.Method == "GET" {
			r.Method = "CHANGED_BY_BEFORE_MIDDLEWARE"
		}
		// w.Write([]byte("Middleware wrote this, so the application should exit.\n"))
	})

	// Adding routin filters, this applies now to "/page/{id}" route
	MiddleWare.Filter("/page/{id}", func(w http.ResponseWriter, r *http.Request, params Frodo.Params) {
		if params.Get("id") != "" {
			r.Method = "CHANGED_BY_FILTER_MIDDLEWARE"
		}
	}, false)

}

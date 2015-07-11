package filters

import (
	"github.com/kn9ts/frodo"
	"net/http"
)

var MiddleWare = Frodo.NewFilters()

func init() {

	MiddleWare.Before(func(w http.ResponseWriter, r *http.Request, _ Frodo.Params) {
		if r.Method == "GET" {
			r.Method = "CHANGED_BY_MIDDLEWARE"
		}
	})

	MiddleWare.Filter("/page/{id}", func(w http.ResponseWriter, r *http.Request, _ Frodo.Params) {
		if r.Method == "GET" {
			r.Method = "CHANGED_BY_MIDDLEWARE"
		}
	}, false)

}

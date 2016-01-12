package frodo

import "net/http"

// MethodsNotAllowed wraps the logic used when an error occurs because the method is not allowed
func MethodsNotAllowed(r *Router, w http.ResponseWriter, req *http.Request) {
	for method := range r.trees {
		// Skip the requested method - we already tried this one
		if method == req.Method {
			continue
		}

		handle, ps, _ := r.trees[method].getValue(req.URL.Path)
		if handle != nil {
			if r.MethodNotAllowed != nil {
				mw := &Middleware{
					Params: ps,
				}
				r.MethodNotAllowed(w, req, mw)
			} else {
				http.Error(w,
					http.StatusText(http.StatusMethodNotAllowed),
					http.StatusMethodNotAllowed,
				)
			}
			return
		}
	}
}

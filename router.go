package frodo

import (
	"fmt"
	"net/http"
)

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	trees map[string]*node

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// Configurable http.Handler which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	NotFound Handle

	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	MethodNotAllowed Handle

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler Handle
}

// Make sure the Router conforms with the http.Handler interface
var _ http.Handler = New()

// New returns a new initialized Router.
// Path auto-correction, including trailing slashes, is enabled by default.
func New() *Router {
	return &Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
	}
}

// Get is a shortcut for router.Handle("GET", path, handle)
func (r *Router) Get(path string, handle ...Handle) {
	r.Handle("GET", path, handle...)
}

// Head is a shortcut for router.Handle("HEAD", path, ...handle)
func (r *Router) Head(path string, handle ...Handle) {
	r.Handle("HEAD", path, handle...)
}

// Options is a shortcut for router.Handle("OPTIONS", path, ...handle)
func (r *Router) Options(path string, handle ...Handle) {
	r.Handle("OPTIONS", path, handle...)
}

// Post is a shortcut for router.Handle("POST", path, ...handle)
func (r *Router) Post(path string, handle ...Handle) {
	r.Handle("POST", path, handle...)
}

// Put is a shortcut for router.Handle("PUT", path, ...handle)
func (r *Router) Put(path string, handle ...Handle) {
	r.Handle("PUT", path, handle...)
}

// Patch is a shortcut for router.Handle("PATCH", path, ...handle)
func (r *Router) Patch(path string, handle ...Handle) {
	r.Handle("PATCH", path, handle...)
}

// Delete is a shortcut for router.Handle("DELETE", path, ...handle)
func (r *Router) Delete(path string, handle ...Handle) {
	r.Handle("DELETE", path, handle...)
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) Handle(method, path string, handle ...Handle) {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if r.trees == nil {
		r.trees = make(map[string]*node)
	}

	root := r.trees[method]
	if root == nil {
		root = new(node)
		r.trees[method] = root
	}

	// Instantiate a Handler type variable
	var handlers []Handle
	// Add all the handlers to it
	handlers = append(handlers, handle...)
	fmt.Printf("%v and the no %d\n", handlers, len(handle))

	// store them to it's route node
	root.addRoute(path, handlers)
}

// Handler is an adapter which allows the usage of an http.Handler as a
// request handle.
func (r *Router) Handler(method, path string, handler http.Handler) {
	r.Handle(method, path,
		func(w http.ResponseWriter, req *http.Request, m *Middleware) {
			handler.ServeHTTP(w, req)
		},
	)
}

// HandlerFunc is an adapter which allows the usage of an http.HandlerFunc as a
// request handle.
func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.Handler(method, path, handler)
}

// ServeFiles serves files from the given file system root.
// The path must end with "/*filepath", files are then served from the local
// path /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use http.Dir:
//     router.ServeFiles("/src/*filepath", http.Dir("/var/www"))
func (r *Router) ServeFiles(path string, root http.FileSystem) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}

	fileServer := http.FileServer(root)

	r.Get(path, func(w http.ResponseWriter, req *http.Request, m *Middleware) {
		req.URL.Path = m.GetParam("filepath")
		fileServer.ServeHTTP(w, req)
	})
}

func (r *Router) recover(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(w, req, nil)
	}
}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handle function and the path parameter
// values. Otherwise the third return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *Router) Lookup(method, path string) ([]Handle, Params, bool) {
	if root := r.trees[method]; root != nil {
		return root.getValue(path)
	}
	return nil, nil, false
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.PanicHandler != nil {
		defer r.recover(w, req)
	}

	if root := r.trees[req.Method]; root != nil {
		path := req.URL.Path

		// get the Handle of the route path requested
		handlers, ps, tsr := root.getValue(path)

		// if []Handle was found were found, run it!
		noOfHandlers := len(handlers)
		if noOfHandlers != 0 {
			// if the 1st handler is defined, run it
			mdwr := &Middleware{
				handlers:     handlers[0:noOfHandlers],
				total:        noOfHandlers,
				nextPosition: 1,
				Params:       ps,
			}
			// run the 1st handler
			// the rest shall be called to run by mdwr.next()
			handlers[0](w, req, mdwr)
		}

		// if a handle was not found, the method is not a CONNECT request
		// and it is not a root path request
		if noOfHandlers == 0 && req.Method != "CONNECT" && path != "/" {
			code := 301 // Permanent redirect, request with GET method
			if req.Method != "GET" {
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308.
				code = 307
			}

			if tsr && r.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}

				http.Redirect(w, req, req.URL.String(), code)
				return
			}

			// Try to fix the request path
			if r.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					CleanPath(path),
					r.RedirectTrailingSlash,
				)
				if found {
					req.URL.Path = string(fixedPath)
					http.Redirect(w, req, req.URL.String(), code)
					return
				}
			}
		}
	}

	// Handle 405
	if r.HandleMethodNotAllowed {
		MethodsNotAllowed(r, w, req)
		return
	}

	//Handle 404
	if r.NotFound != nil {
		r.NotFound(w, req, nil)
	}

	// If system default if CustomHandle not Filter ||
	http.Error(w, http.StatusText(404), http.StatusNotFound)
	return
}

package frodo

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	NotFoundHandler Handle

	// Configurable http.Handler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
	MethodNotAllowedHandler Handle

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler Handle
}

// Make sure the Router conforms with the http.Handler interface
var _ http.Handler = New()

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

// Match adds the Handle to the provided Methods/HTTPVerbs for a given route
// EG. GET/POST from /home to have the same Handle
func (r *Router) Match(httpVerbs Methods, path string, handle ...Handle) {
	if len(httpVerbs) > 0 {
		for _, verb := range httpVerbs {
			r.Handle(strings.ToUpper(verb), path, handle...)
		}
	}
}

// Any method adds the Handle to all HTTP methods/HTTP verbs for the route given
// it does not add routing Handlers for HEADER and OPTIONS HTTP verbs
func (r *Router) Any(path string, handle ...Handle) {
	r.Match(Methods{"GET", "POST", "PUT", "DELETE", "PATCH"}, path, handle...)
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

// Handler is an adapter which allows the usage of an
// http.Handler as a request handle.
func (r *Router) Handler(method, path string, handler http.Handler) {
	r.Handle(method, path, func(w http.ResponseWriter, req *Request) {
		handler.ServeHTTP(w, req.Request)
	})
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

	r.Get(path, func(w http.ResponseWriter, req *Request) {
		req.URL.Path = req.GetParam("filepath")
		fileServer.ServeHTTP(w, req.Request)
	})
}

// NotFound can be used to define custom routes to handle NotFound routes
func (r *Router) NotFound(handler Handle) {
	r.NotFoundHandler = handler
}

// MethodNotAllowed can be used to define custom routes
// to handle Methods that are not allowed
func (r *Router) MethodNotAllowed(handler Handle) {
	r.MethodNotAllowedHandler = handler
}

// ServerError can be used to define custom routes to handle OnServerError routes
func (r *Router) ServerError(handler Handle) {
	r.PanicHandler = handler
}

// On404 is shortform for NotFound
func (r *Router) On404(handler Handle) {
	r.NotFound(handler)
}

// On405 is shortform for NotFound
func (r *Router) On405(handler Handle) {
	r.MethodNotAllowed(handler)
}

// On500 is shortform for ServerError
func (r *Router) On500(handler Handle) {
	r.ServerError(handler)
}

func (r *Router) recover(w *ResponseWriter, req *Request) {
	if err := recover(); err != nil {
		// if a custom panic handler has been defined
		// run that instead
		if r.PanicHandler != nil {
			r.PanicHandler(w, req)
			return
		}

		// If it doesnt, use original http error function as fallback
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
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
	// 1st things 1st, wrap the response writter
	// to add the extra functionality we want basically
	// trace when a write happens
	FrodoWritter := ResponseWriter{
		ResponseWriter: w,
		timeStart:      time.Now(),
		method:         req.Method,
		route:          req.URL.Path,
	}

	// Wrap the supplied http.Request
	FrodoRequest := Request{
		Request: req,
		// params, form - map[string]string,
		// files []*UploadFile
	}

	// ---------- Handle 500: Internal Server Error -----------
	// If a panic/error takes place while process,
	// recover and run PanicHandle if defined
	defer r.recover(&FrodoWritter, &FrodoRequest)

	if root := r.trees[req.Method]; root != nil {
		path := req.URL.Path

		// get the Handle of the route path requested
		handlers, ps, tsr := root.getValue(path)

		// if []Handle was found were found, run it!
		noOfHandlers := len(handlers)
		if noOfHandlers > 0 {
			// if the 1st handler is defined, run it
			FrodoRequest := &Request{
				Request:      req,
				handlers:     handlers[:noOfHandlers-1],
				total:        noOfHandlers,
				nextPosition: 0,
				Params:       ps,
			}

			// call out the middleware Handles
			// the rest shall be called to run by m.Next()
			FrodoRequest.runHandleChain(&FrodoWritter)
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
		for method := range r.trees {
			// Skip the requested method - we already tried this one
			if method == req.Method {
				continue
			}

			handle, ps, _ := r.trees[method].getValue(req.URL.Path)
			if handle != nil {
				if r.MethodNotAllowedHandler != nil {
					FrodoRequest.Params = ps
					r.MethodNotAllowedHandler(&FrodoWritter, &FrodoRequest)
					return
				}
				// if no MethodNotAllowedHandler found, just throw an error the old way
				http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
				return
			}
		}
		return
	}

	// Handle 404
	if r.NotFoundHandler != nil {
		r.NotFoundHandler(&FrodoWritter, &FrodoRequest)
		return
	}

	// If there is not Handle for a 404 error use Go's w
	http.Error(w, http.StatusText(404), http.StatusNotFound)
	return
}

// Serve deploys the application
// Default port is 3102, inspired by https://en.wikipedia.org/wiki/Fourth_Age
// The "Fourth Age" followed the defeat of Sauron and the destruction of his One Ring,
// but did not officially begin until after the Bearers of the Three Rings left Middle-earth for Valinor,
// the 'Uttermost West'
func (r *Router) Serve() {
	r.ServeOnPort(3102)
}

// ServeOnPort is to used if you plan change the port to serve the application on
func (r *Router) ServeOnPort(portNumber interface{}) {
	var portNumberAsString string
	// Converting an interface into the data type it should be
	if pns, ok := portNumber.(int); ok {
		portNumberAsString = strconv.Itoa(pns)
	} else {
		// if it is not a number/int provided then it must be a string
		if pns, ok := portNumber.(string); ok {
			if pns == "" {
				pns = "3102"
			}
			portNumberAsString = pns
		} else {
			log.Fatal("[ERROR] PortNumber can only be a numeral string or integer")
			return
		}
	}

	err := http.ListenAndServe(":"+portNumberAsString, r)
	if err != nil {
		log.Fatalf("[ERROR] Server failed to initialise: %s", err)
		return
	}

	// If server successfully Launched
	log.Printf("[LOG] Server deployed at: %s", portNumberAsString)
}

package frodo

import (
	"fmt"
	"net/http"
	// "reflect"
	"strings"
)

// New global var I will use to launch app
var New *Router

// Handle is a function that can be registered to a route to handle HTTP
// requests. Like http.HandlerFunc, but has a third parameter for the values of
// wildcards (variables).
type Handle func(http.ResponseWriter, *http.Request)
type Handler interface {
	ServerHTTP(http.ResponseWriter, *http.Request)
}

type route struct {
	pattern     string
	parentRoute string
	name        string
	handler     Handle
	isRegex     interface{}
	priority    int
	depth       int
}

type Param struct {
	key   string
	value interface{}
}
type Params []Param

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	paths map[string][]route
}

func NewRouter() *Router {
	New = new(Router)
	return New
}

// Make sure the Router conforms with the http.Handler interface
// var _ http.Handler = NewRouter()
func (r *Router) Application() *Router {
	New = new(Router)
	return New
}

// Get is a shortcut for router.Add("GET", pattern, handle)
func (r *Router) Get(pattern string, handle Handle) {
	r.Handle("GET", pattern, handle)
}

// Head is a shortcut for router.Add("HEAD", pattern, handle)
func (r *Router) Head(pattern string, handle Handle) {
	r.Handle("HEAD", pattern, handle)
}

// Options is a shortcut for router.Add("OPTIONS", pattern, handle)
func (r *Router) Options(pattern string, handle Handle) {
	r.Handle("OPTIONS", pattern, handle)
}

// Post is a shortcut for router.Add("POST", pattern, handle)
func (r *Router) Post(pattern string, handle Handle) {
	r.Handle("POST", pattern, handle)
}

// Put is a shortcut for router.Add("PUT", pattern, handle)
func (r *Router) Put(pattern string, handle Handle) {
	r.Handle("PUT", pattern, handle)
}

// Patch is a shortcut for router.Add("PATCH", pattern, handle)
func (r *Router) Patch(pattern string, handle Handle) {
	r.Handle("PATCH", pattern, handle)
}

// Delete is a shortcut for router.Add("DELETE", pattern, handle)
func (r *Router) Delete(pattern string, handle Handle) {
	r.Handle("DELETE", pattern, handle)
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) Handle(verb, pattern string, handler Handle) bool {
	var routeExists bool
	// word := "UPDATE"
	// capitalise the word if in lowercase
	httpVerb := strings.ToUpper(verb)

	// The route from the Routes Map
	_, exists := r.paths[httpVerb]

	var isReg int
	isReg = strings.IndexAny(pattern, "{}")
	patternInfo := strings.Split(pattern, "/")

	newRoute := route{
		pattern:     pattern,
		parentRoute: patternInfo[0],
		name:        "",
		handler:     handler,
		isRegex:     isReg / 2,
		priority:    0,
		depth:       len(patternInfo),
	}

	// Elem := reflect.ValueOf(handler).Type()
	// fmt.Println(Elem)
	// fmt.Println("Adding these paths to the Router | %s | %s | %b", verb, pattern, exists)

	// If the route map exists r["GET"], r["POST"]...etc`
	if exists {
		// loop thru the list of routes
		for _, rt := range r.paths[httpVerb] {
			// check to see if the route was added before
			if rt.pattern == pattern {
				routeExists = true
			}
		}

		// If has not been added, add it
		if !routeExists {
			r.paths[httpVerb] = append(r.paths[httpVerb], newRoute)
			// fmt.Println(r.paths[httpVerb])
		}
		return true
	} else {
		// fmt.Println("Zero routes added, must initialise and then add")
		// initialise the path map
		r.paths = make(map[string][]route)
		// add the 1st path
		r.paths[httpVerb] = append(r.paths[httpVerb], newRoute)
		return true
	}
	return false
}

// Handler is an adapter which allows the usage of an http.Handler as a
// request handle.
func (r *Router) Handler(method, path string, handler http.Handler) {
	r.Handle(method, path,
		func(w http.ResponseWriter, req *http.Request) {
			handler.ServeHTTP(w, req)
		},
	)
}

// HandlerFunc is an adapter which allows the usage of an http.HandlerFunc as a
// request handle.
func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.Handler(method, path, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("A request came in from %s %s", req.Method, req.URL.Path)
	// h := r.paths[strings.ToUpper(req.Method)]
	// r.Handler(req.Method, h[0].pattern, h[0].handler)
}

func (r *Router) Run() {
	fmt.Println("Server deployed at 3000")
	http.ListenAndServe(":3000", r)
}

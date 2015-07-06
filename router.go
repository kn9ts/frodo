package frodo

import (
	"fmt"
	"net/http"
	"strings"
)

type route struct {
	pattern     string
	parentRoute string
	name        string
	handler     interface{}
	isRegex     interface{}
	priority    int
	depth       int
}

// Router Array Map: {"POST": [r,r,r,r,r,r....], "POST": [r,r,r,r,r,r....]}
type Router struct {
	paths map[string][]route
}

type Param struct {
	key   string
	value interface{}
}

type Params []Param
type Handle func(w http.ResponseWriter, r *http.Request)
type Handler interface {
	ServerHTTP(w http.ResponseWriter, r *http.Request)
}

var New *Router

// Application creates a new instance for Routing
func (r *Router) Application() *Router {
	New = new(Router)
	return New
}

// GET is a shortcut for router.Add("GET", pattern, handle)
func (r *Router) Get(pattern string, handle Handle) {
	r.Add("GET", pattern, handle)
}

// HEAD is a shortcut for router.Add("HEAD", pattern, handle)
func (r *Router) Head(pattern string, handle Handle) {
	r.Add("HEAD", pattern, handle)
}

// OPTIONS is a shortcut for router.Add("OPTIONS", pattern, handle)
func (r *Router) Options(pattern string, handle Handle) {
	r.Add("OPTIONS", pattern, handle)
}

// POST is a shortcut for router.Add("POST", pattern, handle)
func (r *Router) Post(pattern string, handle Handle) {
	r.Add("POST", pattern, handle)
}

// PUT is a shortcut for router.Add("PUT", pattern, handle)
func (r *Router) Put(pattern string, handle Handle) {
	r.Add("PUT", pattern, handle)
}

// PATCH is a shortcut for router.Add("PATCH", pattern, handle)
func (r *Router) Patch(pattern string, handle Handle) {
	r.Add("PATCH", pattern, handle)
}

// DELETE is a shortcut for router.Add("DELETE", pattern, handle)
func (r *Router) Delete(pattern string, handle Handle) {
	r.Add("DELETE", pattern, handle)
}

func (r *Router) Add(verb string, pattern string, handler Handle) bool {
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

	fmt.Println("Adding these paths to the Router | %s | %s | %b", verb, pattern, exists)
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
		}
	} else {
		fmt.Println("Zero routes added, must initialise and then add")
		r.paths = make(map[string][]route)
		r.paths[httpVerb] = append(r.paths[httpVerb], newRoute)
	}
	return false
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "Hello Eugene!") // send data to client side
}

func (r *Router) Serve() {
	fmt.Println("Server deployed at 3000")
	http.ListenAndServe("localhost:3000", nil)
}

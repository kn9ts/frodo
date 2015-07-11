package Frodo

import (
	"fmt"
	"net/http"
	"regexp"
	// "regexp/syntax"
	// "reflect"
	"strconv"
	"strings"
)

// New global var is used to launch the app/routing
var New *Router

// Handle is a function that can be registered to a route to handle HTTP
// requests. Like http.HandlerFunc, but has a third parameter for the values of
// wildcards (variables).
type Handle func(http.ResponseWriter, *http.Request, Params)

type route struct {
	pattern string
	handler Handle
	isRegex int
	depth   int
}

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	paths map[string][]route
	Middleware
}

// Application return a new pointed Router instance
func (r *Router) Application() *Router {
	New = new(Router)
	return New
}

// Get is a shortcut for router.Add("GET", pattern, handle)
func (r *Router) Get(pattern string, handle Handle) {
	r.Handle("GET", pattern, handle)
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

// Head is a shortcut for router.Add("HEAD", pattern, handle)
func (r *Router) Head(pattern string, handle Handle) {
	r.Handle("HEAD", pattern, handle)
}

// Options is a shortcut for router.Add("OPTIONS", pattern, handle)
func (r *Router) Options(pattern string, handle Handle) {
	r.Handle("OPTIONS", pattern, handle)
}

// Handle registers a new request handle with the given path and method.
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
func (r *Router) Handle(verb, pattern string, handler Handle) {
	var routeExists bool

	// If it is "/" <-- root directory
	if li := strings.LastIndex(pattern, "/"); li == 0 && (len(pattern)-1) == li {
		pattern = "/root"
	}

	// word := "UPDATE"
	// capitalise the word if in lowercase
	httpVerb := strings.ToUpper(verb)

	// The route from the Routes Map
	_, exists := r.paths[httpVerb]

	// check to see if it is a regex pattern given from dev
	isReg := len(regexp.MustCompile(`\{[\w.-]{2,}\}`).FindAllString(pattern, -1))
	depth := len(strings.Split(pattern[1:], "/"))

	newRoute := route{
		pattern: pattern,
		handler: handler,
		isRegex: isReg / 2,
		depth:   depth,
	}
	fmt.Printf("Adding this paths to the Router.path[%s] | %v\n", httpVerb, newRoute)

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
	} else {
		// initialise the path map, if nothing had been added
		if len(r.paths) == 0 {
			fmt.Println("Zero routes added, must initialise and then add")
			r.paths = make(map[string][]route)
		}
		// add the 1st path
		r.paths[httpVerb] = append(r.paths[httpVerb], newRoute)
	}
}

// Handler is an adapter which allows the usage of an http.Handler as a
// request handle.
func (r *Router) Handler(method, path string, handler http.Handler) {
	r.Handle(method, path,
		func(w http.ResponseWriter, req *http.Request, p Params) {
			handler.ServeHTTP(w, req)
		},
	)
}

// HandlerFunc is an adapter which allows the usage of an http.HandlerFunc as a
// request handle. Stolen idea from 'httprouter'
func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.Handler(method, path, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Get the URL Path
	requestURL := req.URL.String()

	// If it is "/", root request
	// Convert to equive the request of "/root"
	li := strings.LastIndex(requestURL, "/")
	if li == 0 && (len(requestURL)-1) == li {
		requestURL = "/root"
	}

	// Remove the last trailing slash in user's request
	// E.G /user/eugene/
	if li == len(requestURL)-1 {
		requestURL = requestURL[:len(requestURL)-1]
	}

	// Remove the 1st slash "/", so either have "root/someshit/moreshit"
	requestedURLParts := strings.Split(requestURL[1:], "/")
	fmt.Printf("\n ------- A request came in from [%s] %q --------\n\n", req.Method, req.URL.String())
	fmt.Printf("Requested URL parts -- %q, %d \n", requestedURLParts, len(requestedURLParts))
	// fmt.Println(requestedURLParts[len(requestedURLParts)-1] == "/root")

	// Get the method related with the request in  a []routes array
	// The list of routes related with the method of the requested path
	RelatedRoutes := r.paths[strings.ToUpper(req.Method)]
	fmt.Printf("PATH ROUTES: %v \n", RelatedRoutes)

	// Now loop thru all the routes provided in that Method
	for _, route := range RelatedRoutes {
		// By default on start it is FALSE
		aPossibleRouteMatchFound := false

		// Split and compare the depth of the route requested
		patternSplit := strings.Split(route.pattern[1:], "/")
		fmt.Printf("Request comparisons: %q <<>> %q \n", patternSplit, requestedURLParts)

		//  If the depth match, then they might be a possible match
		if route.depth == len(requestedURLParts) {

			aPossibleRouteMatchFound = true

			// Collect the params in the slice knowing
			// the number of params to expect
			requestParams := make(Params)

			// If a possible match was acquired, step 2:
			// loop thru each part matching them, if one fails then it's a no match
			for index, portion := range requestedURLParts {
				// check to see pattern is something like {param}
				isPattern, _ := regexp.MatchString(`\{[\w.-]{2,}\}`, patternSplit[index])

				// if route part value is actually a regex to match e.g {param}
				if isPattern {
					fmt.Printf("Trying to match the pattern: %q <<>> part: %q", patternSplit[index], portion)

					// Does it pass the {param} match, remove curly brackets
					itMatchesRegexParam, _ := regexp.MatchString(`[\w.-]{2,}`, portion)

					// If it does pass the match test
					if itMatchesRegexParam {
						// Replace all curly brackets with nothing to get the key value
						key := regexp.MustCompile(`(\{|\})`).ReplaceAllString(patternSplit[index], "")
						// Add it to the parameters
						if key != "" {
							requestParams[key] = portion
						}
						aPossibleRouteMatchFound = true
					}

				} else {
					// if there is no regex match,
					// try match them side by side
					if patternSplit[index] != portion {
						aPossibleRouteMatchFound = false
						// If no match here, break the search nothing found
						fmt.Printf("\n-------- BREAK: No match found at all. ---------\n\n")
						break
					}
				}
			}

			// After checking the portions if aPossibleRouteMatchFound remains true,
			// then all the portions matched, thus this handler suffices
			// the match thus run it
			if aPossibleRouteMatchFound {
				fmt.Printf("\nParams found: %q\n", requestParams)

				// Wrap the supplied ResponseWriter
				MiddlewareWriter := &MiddlewareResponseWriter{
					ResponseWriter: w,
				}

				// Get Application "before" Middleware and run them
				if len(r.BeforeMiddleware) > 0 {
					for ix, beforeFilter := range r.BeforeMiddleware {
						// Pass it as the ResponseWriter instead
						beforeFilter(MiddlewareWriter, req, requestParams)
						fmt.Printf("\nBEFORE Middleware No. %d running: Written - %s | Request: - %v \n", ix, req.Method, MiddlewareWriter.written)

						// If there was a write, stop processing
						if MiddlewareWriter.written {
							fmt.Printf("\nEXITING: A write was made by Middleware No. %d | %s \n", ix, req.Method)
							// End the connection
							return
						}
					}
				} else {
					fmt.Printf("\n--- NO middleware: %q ---\n", r.BeforeMiddleware)
				}

				// If there is a middleware that should be implemented to
				// the route, it happens here, before the main dev's handler's
				if len(r.FilterMiddleware) > 0 {
					for _, routeFilter := range r.FilterMiddleware {
						// Try find any route filters that match the route pattern that exist and run them
						if routeFilter.Name == route.pattern {
							routeFilter.Handle(MiddlewareWriter, req, requestParams)
							fmt.Printf("\nROUTE Middleware running: Written - %s | Request: - %v \n", req.Method, MiddlewareWriter.written)
							// If there was a write, stop processing
							if MiddlewareWriter.written {
								fmt.Printf("\nEXITING: A write was made by Middleware Name. %d | %s \n", routeFilter.Name, req.Method)
								// End the connection
								return
							}
							// a match was found, break out
							break
						}
					}

				} else {
					fmt.Printf("\n--- NO middleware: %q ---\n", r.BeforeMiddleware)
				}

				route.handler(w, req, requestParams)
				fmt.Printf("\n-------- EXIT: A Match was made. ---------\n\n")
				break
			}
		}
	}
}

// Serve Deploys the application to the route given
func (r *Router) Serve() {
	http.ListenAndServe(":3000", r)
}

// ServeOnPort is to used if you change the port that you plan on serving on
func (r *Router) ServeOnPort(portNumber interface{}) {
	// Converting an interface into the data type it should be
	portNumber = portNumber.(int)
	if portNumber == 0 {
		portNumber = 3000
	}

	portNumberString := strconv.Itoa(portNumber.(int))
	fmt.Println("Server deployed at: " + portNumberString)
	http.ListenAndServe(":"+portNumberString, r)
}

// AddFilters add Middlewares to routes, requests and responses
func (r *Router) AddFilters(m *Middleware) {
	fmt.Printf("\nMiddleware added %q\n", m.BeforeMiddleware)
	r.Middleware.BeforeMiddleware = m.BeforeMiddleware
	r.Middleware.AfterMiddleware = m.AfterMiddleware
	r.Middleware.FilterMiddleware = m.FilterMiddleware
}

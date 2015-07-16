package Frodo

import (
	"fmt"
	"net/http"
	"regexp"
	// "regexp/syntax"
	"reflect"
	"strconv"
	"strings"
)

// New global var is used to launch the app/routing
var New *Router

// Handle should be able to carry a HandleFunc or a Controller
// thus only both can satisfy an interface
type Handle interface{}

// HandleFunc is a function that can be registered to a route to handle HTTP requests.
// http.Request is now fused with Params, Inputs and Uploads custom handlers
// these are Facades to handle common request processing cases eg. saving a file
type HandleFunc func(http.ResponseWriter, *Request)

// properties of a single route
type route struct {
	pattern string
	handler interface{}
	isRegex int
	depth   int
	Use     // Handles extra details
}

// Methods type is used in Match method to get all methods user wants to apply
// that will help in invoking the related handler
type Methods []string

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

// Get is a shortcut for router.addHandle("GET", args...)
func (r *Router) Get(args ...interface{}) {
	r.addHandle("GET", args...)
}

// Post is a shortcut for router.addHandle("POST", args...)
func (r *Router) Post(args ...interface{}) {
	r.addHandle("POST", args...)
}

// Put is a shortcut for router.addHandle("PUT", args...)
func (r *Router) Put(args ...interface{}) {
	r.addHandle("PUT", args...)
}

// Patch is a shortcut for router.addHandle("PATCH", args...)
func (r *Router) Patch(args ...interface{}) {
	r.addHandle("POST", args...)
}

// Delete is a shortcut for router.addHandle("DELETE", args...)
func (r *Router) Delete(args ...interface{}) {
	r.addHandle("DELETE", args...)
}

// Head is a shortcut for router.addHandle("HEAD", args...)
func (r *Router) Head(args ...interface{}) {
	r.addHandle("HEAD", args...)
}

// Options is a shortcut for router.addHandle("OPTIONS", args...)
func (r *Router) Options(args ...interface{}) {
	r.addHandle("OPTIONS", args...)
}

// Match adds the Handle to the provided Methods/HTTPVerbs for a given route
// EG. GET/POST from /home to have the same Handle
func (r *Router) Match(httpVerbs Methods, args ...interface{}) {
	if len(httpVerbs) > 0 {
		for _, verb := range httpVerbs {
			r.addHandle(strings.ToUpper(verb), args...)
		}
	}
}

// All method adds the Handle to all Methods/HTTPVerbs for a given route
func (r *Router) All(args ...interface{}) {
	methods := Methods{"GET", "POST", "PATCH", "PUT", "DELETE", "PATCH"}
	r.Match(methods, args...)
}

// addHandle registers a new request handle with the given path and method.
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
func (r *Router) addHandle(verb string, args ...interface{}) {

	// at least 2 and a max of 3 arguments are suppost to be provided
	if len(args) > 1 && len(args) < 4 {

		// check if the 1st parameter is a string
		if _, isString := args[0].(string); !isString {
			Log.Error("Error: expected pattern arguement expecting a string")
			return
		}
		pattern := args[0].(string)

		// Check to see if a HandleFunc was provided if not
		v := reflect.ValueOf(args[1]).Type()
		Log.Info("==> %s", v)

		// Debug: First of check if it is a Frodo.HandleFunc type, might have been altered on first/previous loop
		// if not check the function if it suffices the HandleFunc type pattern
		// If it does -- func(http.ResponseWriter, *Request)
		// then convert it to a Frodo.HandleFunc type
		// this becomes neat since this what we expect to run
		// isHandleFunc := false
		if _, ok := args[1].(HandleFunc); !ok {
			if value, ok := args[1].(func(http.ResponseWriter, *Request)); ok && v.Kind().String() == "func" {
				makeHandler := func(h HandleFunc) HandleFunc {
					Log.Debug("converting func(http.ResponseWriter, *Request) to Frodo.HandleFunc")
					return h
				}
				// morph it to it's dynamic data type
				args[1] = makeHandler(value)
				// isHandleFunc = true
			} else {
				// further checked if it is a Controller
				if _, isController := args[1].(ControllerInterface); !isController {
					Log.Fatal("Error: expected handler arguement provided to be an extension of Frodo.Controller or \"func(http.ResponseWriter, *Frodo.Request)\" type")
				}
				args[1] = args[1].(ControllerInterface)
			}
		}
		handler := args[1]
		Log.Info("---- %q ----", reflect.ValueOf(args[1]).Type().String())

		// if the arguments are 3
		isString := false
		if len(args) > 2 {
			// check if the meta/controller information type is Frodo.Use
			if _, isUseStruct := args[2].(Use); !isUseStruct {
				// if it is not Frodo.Use, then check if it is a string
				// probably we were just given the name of the route
				if _, isString := args[2].(string); !isString {
					// If all the tests have passed,
					Log.Fatal("Error: expected controller informative argument provided to be a string or Frodo.Use type")
				}
				args[2] = args[2].(string)
				isString = true
			} else {
				args[2] = args[2].(Use)
			}
			// we now have pattern, handle and info/name
			Log.Debug("pattern: %s | handle: %q | use/name: %q\n", pattern, reflect.ValueOf(args[1]).Type(), reflect.ValueOf(args[2]).Type())
		} else {
			// only pattern and handler are given
			Log.Debug("pattern: %s | handle: %q\n", pattern, reflect.ValueOf(args[1]).Type())
		}

		var routeExists bool

		// If it is "/" <-- root directory
		if li := strings.LastIndex(pattern, "/"); li == 0 && (len(pattern)-1) == li {
			pattern = "/root"
		}

		// word := "GET", "POST", "UPDATE"
		// capitalise the word if it is in lowercase
		httpVerb := strings.ToUpper(verb)

		// Check to see if there is a Routes Map Array for the given HTTP Verb
		_, exists := r.paths[httpVerb]

		// check to see if it is a regex pattern given from dev
		isReg := len(regexp.MustCompile(`\{[\w.-]{2,}\}`).FindAllString(pattern, -1))
		depth := len(strings.Split(pattern[1:], "/"))

		newRoute := route{}
		newRoute.pattern = pattern
		newRoute.handler = handler
		newRoute.isRegex = isReg / 2
		newRoute.depth = depth

		// Add if the name or  meta data of route, if they were given
		if len(args) > 2 {
			if isString {
				newRoute.Name = args[2].(string)
			} else {
				newRoute.Use = args[2].(Use)
			}
		}

		// If the route map exists r["GET"], r["POST"]...etc`
		if exists {
			// loop thru the list of existing routes
			for _, rt := range r.paths[httpVerb] {
				// check to see if the route already exists
				if rt.pattern == pattern {
					routeExists = true
				}
			}

			// If it has not been added, add it
			if !routeExists {
				r.paths[httpVerb] = append(r.paths[httpVerb], newRoute)
			}
		} else {
			// initialise the path map, if nothing had been added
			if len(r.paths) == 0 {
				Log.Warn("Zero routes added, must initialise and then add")
				r.paths = make(map[string][]route)
			}
			// add the 1st path
			r.paths[httpVerb] = append(r.paths[httpVerb], newRoute)
		}

		Log.Success("Adding this route[%v] for the METHOD[%s]\n", httpVerb, newRoute)
	} else {
		// not enough arguements provided
		Log.Fatal("Error: Not enough arguements provided.")
	}
}

// Handler is an adapter which allows the usage of an http.Handler as a
// request handle.
func (r *Router) Handler(method, path string, handler http.Handler) {
	r.addHandle(method, path,
		func(w http.ResponseWriter, req *Request) {
			handler.ServeHTTP(w, req.Request)
		},
	)
}

// HandlerFunc is an adapter which allows the usage of an http.HandlerFunc as a
// request handle. Stolen idea from 'httprouter'
func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.Handler(method, path, handler)
}

// ServeHTTP will receive all requests, and process them for our router
// By using it we are implementing the http.Handler and thus can use our own ways to
// handle incoming requests and process them
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Get the URL Path
	requestURL := req.URL.String()

	// If it is "/", root request
	// Convert it to equive the request of "/root"
	li := strings.LastIndex(requestURL, "/")
	if li == 0 && (len(requestURL)-1) == li {
		requestURL = "/root"
	}

	// Remove the last trailing slash in user's request
	// E.G http://localhost:3000/user/eugene/
	if li == len(requestURL)-1 {
		requestURL = requestURL[:len(requestURL)-1]
	}

	// Remove the 1st slash "/", so either have "root/someshit/moreshit"
	// This helps it matched the stored route paths
	requestedURLParts := strings.Split(requestURL[1:], "/")
	fmt.Printf("\n------- A request came in from [%s] %q --------\n\n", req.Method, req.URL.String())
	fmt.Printf("Requested URL parts -- %q, %d \n", requestedURLParts, len(requestedURLParts))
	// fmt.Println(requestedURLParts[len(requestedURLParts)-1] == "/root")

	// Get list of routes related with the method of the requested path
	// returns a []route array
	RelatedRoutes := r.paths[strings.ToUpper(req.Method)]
	// fmt.Printf("PATH ROUTES: %v \n", RelatedRoutes)

	// Now loop through all the routes provided in that Method
	for _, route := range RelatedRoutes {
		// By default on start it is FALSE
		aPossibleRouteMatchFound := false

		// Split and compare the depth of the route requested
		patternSplit := strings.Split(route.pattern[1:], "/")
		fmt.Printf("Request comparisons: %q <<>> %q \n", patternSplit, requestedURLParts)

		//  If the depth match, then they might be a possible match
		if route.depth == len(requestedURLParts) {
			// For now it seems it is true
			aPossibleRouteMatchFound = true

			// Collect the params in the Params array
			requestParams := &Request{
				Request: req,
				// params, form map[string]string,
				// files []*UploadFile
			}

			// If a possible match was acquired, step 2:
			// loop thru each part of the route pattern matching them, if one fails then it's a no match
			// each part is seperated with "/"
			for index, portion := range requestedURLParts {
				// check to see route part is a pattern eg. {param}
				isPattern, _ := regexp.MatchString(`\{[\w.-]{2,}\}`, patternSplit[index])

				// if the route part value is actually a regex to match
				if isPattern {
					fmt.Printf("Trying to match the pattern: %q <-|-> part: %q", patternSplit[index], portion)

					// Does it pass the {param} match, remove curly brackets
					routePartMatchesRegexParam, _ := regexp.MatchString(`[\w.-]{2,}`, portion)

					// If it does pass the match test
					if routePartMatchesRegexParam {
						// Replace all curly brackets with nothing to get the key value
						key := regexp.MustCompile(`(\{|\})`).ReplaceAllString(patternSplit[index], "")
						// Add it to the parameters
						if requestParams.params == nil {
							requestParams.params = make(map[string]string)
						}

						if key != "" {
							requestParams.params[key] = portion
						}
						// Keep it true
						aPossibleRouteMatchFound = true
					}

				} else {
					// if there is no regex match, try match them side by side as strings
					// If no match was found, then we are wasting time
					if patternSplit[index] != portion {
						// If no match here, falsify & break the search nothing found
						aPossibleRouteMatchFound = false
						fmt.Printf("\n-------- BREAK: No match found at all. ---------\n\n")
						break
					}
				}
			}

			// After checking the portions if aPossibleRouteMatchFound remains true,
			// then all the portions matched, thus this route suffices the match
			// thus grab it's handler and run it
			if aPossibleRouteMatchFound {
				// fmt.Printf("\nParams found: %q\n", requestParams)

				// Wrap the supplied http.ResponseWriter, we want to know when
				// a write has been done by the middleware or controller and exit immediately
				MiddlewareWriter := &MiddlewareResponseWriter{
					// Since http.ResponseWriter is embedded you can access it
					ResponseWriter: w,
				}

				// Get Application's "Before" Middleware and run them
				if len(r.BeforeMiddleware) > 0 {
					for ix, _ := range r.BeforeMiddleware {
						// Pass it as the ResponseWriter instead
						// beforeFilter(MiddlewareWriter, requestParams)
						fmt.Printf("\nBEFORE Middleware No. %d running: Written - %s | Request: - %v \n", ix, req.Method, MiddlewareWriter.written)

						// If there was a write, stop processing
						if MiddlewareWriter.written {
							fmt.Printf("\nEXITING: A write was made by Middleware Number: %d | %s \n", ix, req.Method)
							// End the connection
							return
						}
					}
				} else {
					// No before middleware added
					fmt.Printf("\n--- NO middleware: %q ---\n", r.BeforeMiddleware)
				}

				// If there is a middleware that should be implemented to the route
				// Run it before the controller, before the controller provided by the dev
				if len(r.FilterMiddleware) > 0 {
					// loop thru the filters to find it
					for _, routeFilter := range r.FilterMiddleware {
						// Try find any route filter that matches the route pattern and run them
						if routeFilter.Name == route.pattern {
							// routeFilter.Handle(MiddlewareWriter, requestParams)
							fmt.Printf("\nROUTE Middleware running: Written - %s | Request: - %v \n", req.Method, MiddlewareWriter.written)
							// If there was a write, stop processing
							if MiddlewareWriter.written {
								fmt.Printf("\nEXITING: A write was made by Middleware Name: %d | %s \n", routeFilter.Name, req.Method)
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

				// Finally run the dev's controller provided, and exit (for now though, after middleware to be added)

				/*

				   ------- RESULTS: -----

				   Value:  <*Frodo.Controller Value>
				   Type:  *Frodo.Controller
				   Kind:  ptr
				   Interface:  &{Get }
				   Pointer:  833357997216
				   Elem:  <Frodo.Controller Value>

				*/
				v := reflect.ValueOf(route.handler)
				fmt.Println("Type: ", v.Type().String())
				// route.handler(w, requestParams)

				fmt.Printf("\n-------- EXIT: A Match was made. ---------\n\n")
				if MiddlewareWriter.written {
					// End the connection
					return
				}
				break
			}
		}
	}
}

// Serve deploys the application
// Default port is 3102, inspired by https://en.wikipedia.org/wiki/Fourth_Age
// The Fourth Age followed the defeat of Sauron and the destruction of his One Ring,
// but did not officially begin until after the Bearers of the Three Rings left Middle-earth for Valinor,
// the 'Uttermost West'
func (r *Router) Serve() {
	r.ServeOnPort(3102)
}

// ServeOnPort is to used if you change the port that you plan on serving on
func (r *Router) ServeOnPort(portNumber interface{}) {
	var portNumberString string
	// Converting an interface into the data type it should be
	if pn, ok := portNumber.(int); ok {
		portNumberString = strconv.Itoa(pn)
	} else {
		// if it is not a number/int provided then it must be a string
		if pns, ok := portNumber.(string); ok {
			if pns == "" {
				pns = "3102"
			}
			portNumberString = pns
		} else {
			Log.Fatal("Error: PortNumber can only be a numeral string or integer")
		}
	}

	err := http.ListenAndServe(":"+portNumberString, r)
	if err != nil {
		Log.Fatal("Error: server failed to initialise: %s", err)
	}
	// If server successfully Launched
	Log.Success("Server deployed at: %s", portNumberString)
}

// AddFilters add Middlewares to routes, requests and responses
func (r *Router) AddFilters(m *Middleware) {
	fmt.Printf("\nMiddleware added %q\n", m.BeforeMiddleware)
	r.Middleware.BeforeMiddleware = m.BeforeMiddleware
	r.Middleware.AfterMiddleware = m.AfterMiddleware
	r.Middleware.FilterMiddleware = m.FilterMiddleware
}

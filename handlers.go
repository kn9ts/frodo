package frodo

import (
	"net/http"
)

// Handle should be able to carry a HandleFunc or a Controller
// thus only both can satisfy an interface
// type Handle interface{}

// Handle is a function that can be registered to a route to handle HTTP
// requests. Like http.HandlerFunc, but has a third parameter for the values of
// wildcards (variables).
type Handle func(http.ResponseWriter, *http.Request, Params)

// HandleFunc is a function that can be registered to a route to handle HTTP requests.
// http.Request is now fused with Params, Inputs and Uploads custom handlers
// these are Facades to handle common request processing cases eg. saving a file
type HandleFunc func(w http.ResponseWriter, r *http.Request)

// HTTPRouterHandleFunc describles the httprouter handler
type HTTPRouterHandleFunc Handle

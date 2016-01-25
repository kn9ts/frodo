package frodo

import (
	"net/http"
)

// Controller defines the basic structure of a REST application controller,
// the devs controller should embed this to create their own controllers.
// It then automatically implements ControllerInterface, and can be passed as a controller in routing
type Controller struct {
	Method string
	Attributes
}

// ControllerInterface will be used to parse back the developer's controller back to it's own type.
// Since we know that a REST controller entails the following methods then any struct that implements the Controller methods suffices the ControllerInterface
type ControllerInterface interface {
	Index(http.ResponseWriter, *Request)
	Create(http.ResponseWriter, *Request)
	Store(http.ResponseWriter, *Request)
	Show(http.ResponseWriter, *Request)
	Edit(http.ResponseWriter, *Request)
	Update(http.ResponseWriter, *Request)
	Patch(http.ResponseWriter, *Request)
	Destroy(http.ResponseWriter, *Request)
	Head(http.ResponseWriter, *Request)
	Options(http.ResponseWriter, *Request)
	Next(...interface{})
}

// Index is the default handler for any incoming request or route's request that is not matched to it's handler
// It can also be used for specific route, mostly for the root routes("/")
func (c *Controller) Index(w http.ResponseWriter, r *Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Create handles a request to Show the form for creating a new resource.
// * GET /posts/create
func (c *Controller) Create(w http.ResponseWriter, r *Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Store handles a request to Store a newly created resource in storage.
// * POST /posts
func (c *Controller) Store(w http.ResponseWriter, r *Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Show handles a request to Display the specified resource.
// * GET /posts/{id}
func (c *Controller) Show(w http.ResponseWriter, r *Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Edit handles a request to Show the form for editing the specified resource.
// * GET /posts/{id}/edit
func (c *Controller) Edit(w http.ResponseWriter, r *Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Update handles a request to update the specified resource in storage.
// * PUT /posts/{id}
func (c *Controller) Update(w http.ResponseWriter, r *Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Patch is an alternative to Update
func (c *Controller) Patch(w http.ResponseWriter, r *Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Destroy handles a request to Remove the specified resource from storage.
// * DELETE /posts/{id}
func (c *Controller) Destroy(w http.ResponseWriter, r *Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Head handle HEAD request.
func (c *Controller) Head(w http.ResponseWriter, r *Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Options handle OPTIONS request.
func (c *Controller) Options(w http.ResponseWriter, r *Request) {
	http.Error(w, "Method Not Allowed", 405)
}

// Next will be used to call the next handler in line/queue
// th biggest change is that it requires the Request struct to be passed as parametre
// to call the next handler in line
//
// it also makes the Controller implement the Middleware type
func (c *Controller) Next(args ...interface{}) {
	// We only need the 1st argument
	// the Request Object
	// r := args[0].(*Request)
}

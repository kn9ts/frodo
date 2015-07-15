package Frodo

import (
	"net/http"
)

// Controller Defines the basic structure of a REST Application Controller
// the devs, controller should embedd this
type Controller struct {
	Method, Layout string
}

// ControllerInterface will be used to parse back the dev's Controller back to it's own type
// Since we know what a REST controller entails the following methods
// Then any Struct that implements the Controller methods
// suffices the ControllerInterface
type ControllerInterface interface {
	_Construct(func(*Controller))
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
	NotFound(http.ResponseWriter, *Request)
}

func (c *Controller) _Construct(fn func(*Controller)) {
	fn(c)
}

// Index is the default handler for any incoming request that is not matched to it's handler
// It can also be used for specific route, mostly for root routes
// Display a listing of the resource.
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

// NotFound can be customised to be used as an alternative thus run your own custom handles
// eg. Redirect the user to a 404 page, or rooot page
func (c *Controller) NotFound(w http.ResponseWriter, r *Request) {
	http.Redirect(w, r.Request, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

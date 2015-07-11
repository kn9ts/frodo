package Frodo

// CustomHandle is for filter Middleware, belong to a method or route
// Reminder:
// type Handle func(http.ResponseWriter, *http.Request, Params)
type CustomHandle struct {
	Name   string
	Handle Handle
	Method interface{} // should be string
}

// Middleware refers to a chain of handlers that wrap around the web app
// Adding extra functionality
// eg. authenticate user, redirect guest visits, compress response content, rate limit requests, log application events
// This struct splits them into Before, After or Filter(route specific) middlware for requests is received
type Middleware struct {
	BeforeMiddleware, AfterMiddleware []Handle
	FilterMiddleware                  []CustomHandle
}

// NewFilters creates a new *Middleware instance
func NewFilters() *Middleware {
	return new(Middleware)
}

// Before adds Middleware that run as soon as request comes in
func (m *Middleware) Before(h Handle) {
	m.BeforeMiddleware = append(m.BeforeMiddleware, h)
}

// After adds Middleware that run after the controller specified before giving back a response
func (m *Middleware) After(h Handle) {
	m.AfterMiddleware = append(m.AfterMiddleware, h)
}

// Filter adds Middleware that will be special to an application route
func (m *Middleware) Filter(name string, h Handle, method interface{}) {
	m.FilterMiddleware = append(m.FilterMiddleware, CustomHandle{
		Name:   name,
		Handle: h,
		Method: method,
	})
}

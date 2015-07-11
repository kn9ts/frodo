package Frodo

// CustomFilter is for filter Middleware, belong to a method or route
type CustomFilter struct {
	Name   string
	Handle Handle
	Method interface{}
}

// Middleware is plainly an array of functions that run
// before and after once a request is received, even a user's controller is middleware
type Middleware struct {
	BeforeMiddleware, AfterMiddleware []Handle
	FilterMiddleware                  []CustomFilter
}

// NewFilters creates a new *Middleware instance
func NewFilters() *Middleware {
	return new(Middleware)
}

// Before adds Middleware that run as soon as request comes in
func (m *Middleware) Before(h Handle) {
	m.BeforeMiddleware = append(m.BeforeMiddleware, h)
}

// After adds Middleware that run after user's processes and before
// giving back a response in however form it will be sent out
func (m *Middleware) After(h Handle) {
	m.AfterMiddleware = append(m.AfterMiddleware, h)
}

// Filter adds middleware that will be special to a route or event
// EG. Authenticationg the user
func (m *Middleware) Filter(name string, h Handle, method interface{}) {
	m.FilterMiddleware = append(m.FilterMiddleware, CustomFilter{
		Name:   name,
		Handle: h,
		Method: method,
	})
}

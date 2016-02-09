package frodo

// Attributes will be used as an optional argument while declaring routes
// it lets you:
//   - Name a Controller or Handler
//   - define the specific Method to be used in a Controller
//   - a list of Middlewares that should run before the specific Controller
type Attributes struct {
	Method, Name string
	Middleware   []string
}

package Frodo

// Use will be used in conjuction with routing for further sugar in the Routing system
// It defines the Method to be used if the Handle is a Controller
// Filter[string, name of filter] can also be related to the route and run in this specified route only in every incoming requests
type Use struct {
	Method, Name string
	Filter       string
	Meta         map[string]interface{}
}

package Frodo

// Use will be used in conjuction with routing for further sugar in the Routing system
// It defines the Method tob be used if the Handle is a HandleFunc
// and also giving the route a name, that can be used in filters instead of the route.pattern
type Use struct {
	Method, Name  string
	Before, After string
	Meta          map[string]interface{}
}

// func (u *Use) Get(name string) {
//
// }

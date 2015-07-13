package Frodo

import "net/http"

// Request type will carry all the values in curly {} brackets that are
// translated from url param values to ready to be used values
type Request struct {
	*http.Request
	params, form, file map[string]string
}

// Get returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (pm *Request) Get(name string) string {
	value, ok := pm.params[name]
	if ok {
		return value
	}
	return ""
}

// Set adds a key/value pair to the Request params
func (pm *Request) Set(name, value string) bool {
	// 1st check if it has been initialised
	if pm.params != nil { // If not initialise
		pm.params = make(map[string]string)
	}

	// allow overwriting
	pm.params[name] = value
	return true
}

// Input gets ALL posted key/values from all Methods
func (pm *Request) Input() {
}

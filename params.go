package frodo

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params map[string]string

// GetParam returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (p Params) GetParam(name string) string {
	value, ok := p[name]
	if ok {
		return value
	}
	return ""
}

// Param is shorter equivalent of `GetParam` method
func (p Params) Param(name string) string {
	return p.GetParam(name)
}

// SetParam adds a key/value pair to the Request params
func (p Params) SetParam(name, value string) bool {
	// 1st check if it has been initialised
	if p != nil {
		// If not initialise
		p = make(map[string]string)
	}

	// allow overwriting
	p[name] = value
	return true
}

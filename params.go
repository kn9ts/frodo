package Frodo

// Params is passed all routing/url parameters
type Params map[string]string

// Get returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) Get(name string) string {
	value, ok := ps[name]
	if ok {
		return value
	}
	return ""
}

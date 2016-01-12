package frodo

import "net/http"

// Middleware will help facilitate the passing of multiple handlers
type Middleware struct {
	handlers     []Handle
	total        int
	nextPosition int
	Params
}

// Next will be used to call the next handler in line/queue
func (m *Middleware) Next(w http.ResponseWriter, r *http.Request) {
	// 1st check if the next handler position accounts for the number
	// of handlers existing in the handlers array
	if m.nextPosition < m.total {
		// get the next handler
		nextHandler := m.handlers[m.nextPosition]
		// move the cursor
		m.nextPosition++
		// now run the next handler
		nextHandler(w, r, m)
	}
	return
}

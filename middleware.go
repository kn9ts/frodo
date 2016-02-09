package frodo

import "fmt"

// Middleware declares the minimum implementation
// necessary for a handlers to be used as Frodo's middleware route Handlers
type Middleware interface {
	Next(...interface{})
}

// RequestMiddleware is a collection of Middlwares
// this struct will have the method next to invoke the middleware in the chain
// one after the other
type RequestMiddleware struct {
	handlers       []Middleware // []Handle
	ResponseWriter *ResponseWriter
	Request        *Request
	total          int
	nextPosition   int
}

func (m *RequestMiddleware) chainReaction() {
	m.nextPosition++
	m.typeCastAndCall(m.handlers[0])
}

// Next will be used to call the next handler in line/queue
func (m *RequestMiddleware) Next(args ...interface{}) {
	// 1st check if the next handler position accounts for the number
	// of handlers existing in the handlers array
	if m.nextPosition < m.total {
		// get the next handler
		nextHandle := m.handlers[m.nextPosition]
		// move the cursor
		m.nextPosition++
		// 1st check if a write has happened
		// have headers been sent out?
		// meaning a response has been issued out to the client
		// if not call the next handler in line
		if !m.ResponseWriter.HeaderWritten() {
			m.typeCastAndCall(nextHandle)
		}
	}
}

// typeCastAndCall converts the middleware to it's rightful type then calls it
func (m *RequestMiddleware) typeCastAndCall(run Middleware) {
	// 1st check if the route handler is HandleFunc
	if handle, hasTypeCasted := run.(Handler); hasTypeCasted {
		handle(m.ResponseWriter, m.Request)
	} else {
		// if not, then is it an implementation of ControllerInterface
		if ctrl, ok := run.(CRUDController); ok {
			// Yes! it is.
			// Ok! check if a method was specified to run
			// if name := ctrl.Method; name != "" {
			// 	// if so check that Method exists
			// 	v := reflect.ValueOf(ctrl)
			//
			// 	// check for the method by it's name
			// 	fn := v.MethodByName(name)
			//
			// 	// if a Method was found, not a Zero value
			// 	if fn != (reflect.Value{}) {
			// 		// Then convert it back to a Handler
			// 		// You have to know which type it is you are converting to
			// 		if value, ok := fn.Interface().(func(http.ResponseWriter, *Request)); ok && fn.Kind().String() == "func" {
			// 			// morph it to it's dynamic data type, and run it
			// 			makeHandler(value)(m.ResponseWriter, m.Request)
			// 			return
			// 		}
			// 	} else {
			// 		// Method given in use does not exist
			// 		err := fmt.Errorf("Error: Method undefined (The Controller has no field or method %s)", name)
			// 		panic(err)
			// 	}
			// } else {
			// 	// Nothing like so were found, run internal server error: 500
			// 	fmt.Println("No Method specified to run in Controller, defaulting to Index method")
			// 	ctrl.Index(m.ResponseWriter, m.Request)
			// 	return
			// }

			// If no Controller.Attribute.Method was provided, run Index as the default fallback
			ctrl.Index(m.ResponseWriter, m.Request)
		} else {
			// No Handler or Controller was found, run internal server error: 500
			m.ResponseWriter.WriteHeader(404)
			fmt.Fprintf(m.ResponseWriter, "No Frodo.Handle or Frodo.Controller exists to handle the route.")
		}
	}
}

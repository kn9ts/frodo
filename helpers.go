package frodo

import "fmt"

// Little function to convert type "func(http.ResponseWriter, *Request)" to Frodo.HandleFunc
func makeHandler(h Handler) Handler {
	fmt.Println("converting func(http.ResponseWriter, *Request) to Frodo.Handler")
	return h
}

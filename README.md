# Frodo [Go Web Mini Framework]

Frodo is a Go mini web framework inspired by the sweet/beautiful parts that make up
Laravel(PHP), Slim (PHP) and ExpressJS(NodeJS).

I built it to so as to learn Go, and also how frameworks work.

Hello world example:

```go
package main

import (
	"net/http"
	"github.com/kn9ts/frodo"
)

func main()  {
	// Create a new instance of Frodo
	App := Frodo.New()

	// Add the root route
	App.Get("/", func(w http.ResponseWriter, r *Frodo.Request) {
		w.Write([]byte("Hello World! I am Frodo."))
	})

	App.Serve() // Open in browser http://localhost:3102/
}
```

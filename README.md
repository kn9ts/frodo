# Frodo (A Tiny Go Web Framework)

[Frodo](http://godoc.org/github.com/kn9ts/frodo) is a Go micro
web framework inspired by ExpressJS.

<!-- I built it to so as to learn Go, and also how
frameworks work. A big thanks to TJ Holowaychuk too
for the inspiration -->

Are you looking for **[GoDocs Documentation](http://godoc.org/github.com/kn9ts/frodo)**

#### Updates

- Intergrated and using [httprouter](https://github.com/julienschmidt/httprouter)
- Accepts middleware now by default, one or more

#### `Hello world` example

```go
package main

import (
		"net/http"
		"github.com/kn9ts/frodo"
)

func one(w http.ResponseWriter, r *frodo.Request) {
	fmt.Println("Hello, am the 1st middleware!")
	// fmt.Fprint(w, "Hello, I'm 1st!\n")
	r.Next()
}

func two(w http.ResponseWriter, r *frodo.Request) {
	fmt.Println("Hello, am function no. 2!")
	// fmt.Fprint(w, "Hello, am function no. 2!\n")
	r.Next()
}

func three(w http.ResponseWriter, r *frodo.Request) {
	fmt.Println("Hello, am function no 3!")
	fmt.Fprint(w, "Hey, am function no. 3!\n")
}

func nameFunction(w http.ResponseWriter, r *frodo.Request) {
	fmt.Println("Hello there, ", r.GetParam("name"))
	fmt.Fprintf(w, "Hello there, %s!\n", r.GetParam("name"))
}

func main() {
	app := frodo.New()

	app.Get("/", one, two, three)
	app.Get("/hello/:name", one, nameFunction)

	app.Serve()
}
```

#### Coming soon

- Controllers
- Ability to detect CRUB requests and run the right controller method

## Release History

**Version: 0.10.0**

## License

Copyright (c) 2014 **Eugene Mutai**
Licensed under the [MIT license](http://mit-license.org/)

# Frodo (A Tiny Go Web Framework)

[Frodo](http://godoc.org/github.com/kn9ts/frodo) is a Go micro web framework inspired by ExpressJS.

__NOTE:__ _I built it to so as to learn Go, and also how frameworks work. A big thanks to TJ Holowaychuk too
for the inspiration_

Are you looking for the **[GoDocs Documentation](http://godoc.org/github.com/kn9ts/frodo)**

#### Updates

- Intergrated(actually interweaved into the code base) and using [httprouter](https://github.com/julienschmidt/httprouter) as the framework's routing system
- Accepts handlers as middleware now by default, one or more

#### "Hello world" example

The `main.go` file:

```go
package main

import (
		"net/http"
		"github.com/kn9ts/frodo"
)

func main() {
	app := frodo.New()

	app.Get("/", one, two, three)
	app.Get("/hello/:name", one, nameFunction)

	app.Serve()
}
```

And the functions passed as middleware would look like:

```go
package main

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
```

#### To do (after Go sabitcal is over)

- Controllers (which will implement a BaseController)
- Controllers can be mixed with the common handlers as middleware
- Ability to detect CRUD requests and run the right controller method, if a controllers are passed as middleware

## Release History

**Version: 0.10.0**

## License

Copyright (c) 2014 **Eugene Mutai**
Licensed under the [MIT license](http://mit-license.org/)

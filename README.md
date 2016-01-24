# Frodo (Go Web Micro Framework)

[Frodo](http://godoc.org/github.com/kn9ts/frodo) is a Go micro
web framework inspired by ExpressJS.

<!-- I built it to so as to learn Go, and also how
frameworks work. A big thanks to TJ Holowaychuk too
for the inspiration -->

Are you looking for **[GoDocs Documentation](http://godoc.org/github.com/kn9ts/frodo)**

##### `Hello world` example

```go
package main

import (
		"net/http"
		"github.com/kn9ts/frodo"
)

func main()  {
		// Create a new instance of frodo
		App := frodo.New()

		// Add the root route
		App.Get("/", func(w http.ResponseWriter, r *frodo.Request) {
				w.Write([]byte("Hello World!"))
		})

		App.Serve() // Open in browser http://localhost:3102/
}
```

##### A bigger `Hello world` example

```go
package main

import (
	"net/http"
	"github.com/kn9ts/frodo"
	"gopkg.in/unrolled/render.v1"
)

func main()  {
	// Create a new instance of frodo
	App := frodo.New()

	// Yes, you can use the famous old render package to ender your views
	Response := render.New(render.Options{})

	// Add the root route
	App.Get("/", func(w http.ResponseWriter, r *frodo.Request) {
		// if you don't get it
		// --> https://www.youtube.com/watch?v=vjW8wmF5VWc
		w.Write([]byte("Hey, Watch Me Nae Nae!!!"))
	})

	// ------ Controller Awesomeness! ------
	// Passing a controller instead of a callback function, runs Index method by default
	App.Get("/home", &controller.Home{})

	// ----- Methods and Dynamic routes -----
	// You can declare which method in a controller should be called for the specified route
	// Oh yeah! you can name your routes eg. user-profile
	App.Post("/profile/{name}", &controller.Home{}, frodo.Options{
		Method: "Profile",
		Name: "user-profile",
	})

	// ----- Multiple Methods -----
	// How about declaring more than one method to accept a specific Request, HELL YES!!!
	App.Match(frodo.Methods{"GET", "POST"}, "/home", func(w http.ResponseWriter, r *frodo.Request) {
		Response.HTML(w, http.StatusOK, "Hello! We are home!", nil)
	})

	App.Serve() // Open in browser http://localhost:3102/
}
```

## Controllers
From the above example you can observe that **frodo** can also use
`controllers` instead of the usual route handlers.

The `controller` used above with the route mapping would look as described below; placed in
the `controllers` folder, which should be in root directory of your project.

The file name does not matter but the package name matters.
It then should embed `frodo.Controller` struct so as to
inherit all **frodo's** controller functionality.

`controller.Home{}` in `./controllers/home.go` would look like this:

```go
package controller

import (
	"github.com/kn9ts/frodo"
	"net/http"
)

// Home is plays an example of a controller
type Home struct {
	*frodo.Controller
}

// Index is the default route method for "/" route
// also if a route controller method is not found,
// it falls back to the Index method
func (h *Home) Index(w http.ResponseWriter, r *frodo.Request) {
	w.Write([]byte("Hello world, a message from Home controller."))
}

func (h *Home) Profile(w http.ResponseWriter, r *frodo.Request) {
	w.Write([]byte("Hey, Watch Me, " + r.Param("name") + ", Dougie from home controller."))
}
```


## Middleware/Application Filters

**Daaaahh!** Ofcourse there are `MiddleWares` in **frodo**.
You can create a folder named `filter` in your project's folder
and declare your MiddleWare there.

Example: `filters.go` inside the `./filters` folder.

```go
```

## Release History

**Version: 0.9.2 Preview**

## License

Copyright (c) 2014 **Eugene Mutai**
Licensed under the [MIT license](http://mit-license.org/)

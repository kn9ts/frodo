# Frodo [Go Web Micro Framework]

Frodo is a Go mini web framework inspired by the sweet/beautiful parts that make up Laravel(php), Slim (php) and ExpressJS(node.js).
<!-- I built it to so as to learn Go, and also how frameworks work. -->

##### `Hello World` example:

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
        // https://www.youtube.com/watch?v=vjW8wmF5VWc
        w.Write([]byte("Hello World!!!"))
    })

    App.Serve() // Open in browser http://localhost:3102/
}
```

##### A big `Hello world` example:

```go
package main

import (
	"net/http"
	"github.com/kn9ts/frodo"
	"github.com/username/project/filters"
	"github.com/username/project/controller"
	"gopkg.in/unrolled/render.v1"
)

func main()  {
	// Create a new instance of Frodo
	App := Frodo.New()

	// Yes, you can use the famous old render package to ender your views
	Response := render.New(render.Options{})

	// Add the root route
	App.Get("/", func(w http.ResponseWriter, r *Frodo.Request) {
		// https://www.youtube.com/watch?v=vjW8wmF5VWc
		w.Write([]byte("Hey, Watch Me Nae Nae!!!"))
	})

	// ------ Controller Awesomeness! ------
	// Passing a controller instead of a callback function, runs Index method by default
	App.Get("/home", &controller.Home{})

	// ----- Methods and Dynamic routes -----
	// You can declare which method in a controller should be called for the specified route
	// Oh yeah! you can name your routes eg. user-profile
	App.Post("/profile/{name}", &controller.Home{}, Frodo.Use{
		Method: "Profile",
		Name: "user-profile",
	})

	// ----- Multiple Methods -----
	// How about declaring more than one method to accept a specific Request, HELL YES!!!
	App.Match(Frodo.Methods{"GET", "POST"}, "/home", func(w http.ResponseWriter, r *Frodo.Request) {
		Response.HTML(w, http.StatusOK, "hello", nil)
	})

	App.AddFilters(filters.MiddleWare)
	App.Serve() // Open in browser http://localhost:3102/
}
```

## Controllers
From the above example you can observe that **Frodo** can also accept `controllers` instead of the usual callback function passed. The controller used above route mapping would look as follows, placed in the `controllers` folder, the file name does not matter but the package name matters. It then should embed `Frodo.Controller` struct so as to inherit controller functionality.

Example of `controller.Home{}` used above would be:

```go
package controller

import (
	"github.com/kn9ts/frodo"
	"net/http"
)

// Home is plays an example of a controller
type Home struct {
	Frodo.Controller
}

// Index is the default route handler for "/" route
func (h *Home) Index(w http.ResponseWriter, r *Frodo.Request) {
	w.Write([]byte("Hello world, a message from Home controller."))
}

func (h *Home) Profile(w http.ResponseWriter, r *Frodo.Request) {
	w.Write([]byte("Hey, Watch Me, " + r.Param("name") + ", Dougie from home controller."))
}
```


## Middleware/Application Filters
**Owh Yeah!** Ofcourse you saw `Filters or MiddleWare` added just before we initialized the server. So you can create a folder named `filter` and declare your MiddleWare there, for example the above `filter.MiddleWare` would look like this when declared in a file inside the `filters` folder.

```go
package filters

import (
	"github.com/kn9ts/frodo"
	"net/http"
)

// This is how a dev makes a Filter instance
var MiddleWare = Frodo.NewFilters()

func init() {

	// Adding application `Before` filters
	MiddleWare.Before(func(w http.ResponseWriter, r *Frodo.Request) {
		// Do something here
	})

	// Adding application `After` filters
	MiddleWare.After(func(w http.ResponseWriter, r *Frodo.Request) {
		// Do something here
	})

	// Adding route secific filters/middleware
	// pass the name of route pattern of the route you want to run the following middleware
	MiddleWare.Filter("user-profile", func(w http.ResponseWriter, r *Frodo.Request) {
		// do something here before the router named 'user-profile'
		// for example check if the name is 'Eugene'
		if r.Param("name") == "Eugene" {
			// Yes! the name passed in the params is Eugene, it has passed
		} else {
			// Name does not match. Do something if it has not passed
		}
	})

}
```

## Release History
__Version: 0.9.1 Preview__

## License
Copyright (c) 2014 __Eugene Mutai__
Licensed under the [MIT license](http://mit-license.org/)

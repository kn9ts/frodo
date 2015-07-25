package main

import (
	"log"
	"strings"
)

type router struct {
	name, pattern string
	handler       interface{}
	// if {param}, what name was it assigned
	param string
	// method to run, and filter to run before any method
	method, filter string
	isRegex, depth int
}

type path struct {
	router
	node map[rune]map[string]*path
}

var Router *path

func newRouter() *path {
	log.SetFlags(log.Lshortfile)
	return &path{
		router: router{
			name:    "root",
			pattern: "/",
			handler: "root hander for the root or default (\"/\") route",
		},
	}
}

func getCharCodeAt(s string, n int) rune {
	for i, r := range s {
		if i == n {
			return r
		}
	}
	return 0
}

func main() {
	Router := newRouter()
	addNodeTest(Router)
	checkForNodesTest(Router)
}

func addNodeTest(r *path) {
	// Create a node
	np := &path{
		router: router{
			name:    "user-image-id",
			pattern: "/user/{param}/images/{id}",
			handler: "/{id} handler for route /user/{param}/images",
		},
		node: make(map[rune]map[string]*path),
	}

	np1 := &path{
		router: router{
			name:    "user-image-id",
			pattern: "/user/{param}",
			handler: "/{param} handler for route /user/",
		},
		node: make(map[rune]map[string]*path),
	}

	// Try to add the node to thr route
	_, ok := r.addNode(np)

	// Test if you will find the route
	pattern := "/user/{param}/images/{id}"
	p, lp := r.getNode(pattern)
	log.Printf("The NODE: %v | lvl: %d \n", p.router.handler, lp)
	log.Printf("<==============================| the node has been added |================================>\n\n%q | %v\n\n", r, ok)

	_, ok1 := r.addNode(np1)

	pattern1 := "/user/{param}"
	p1, lp1 := r.getNode(pattern1)
	log.Printf("The NODE: %v | lvl: %d \n", p1.router.handler, lp1)
	log.Printf("<==============================| the node has been added |================================>\n\n%q | %v\n\n", r, ok1)
}

func checkForNodesTest(r *path) {
	h1 := router{
		name:    "user",
		pattern: "/user/{name}",
		handler: "/user handler here",
	}

	h2 := router{
		name:    "{name}",
		pattern: "/user/{name}",
		handler: "/{name} handler for the route /user",
	}

	h3 := router{
		name:    "images",
		pattern: "/user/{name}/images",
		handler: "/images handler for the route /user/{name}",
	}

	h4 := router{
		name:    "image-id",
		pattern: "/user/{name}/images/{id}",
		handler: "/{id} handler for the route /user/{name}/images",
	}

	h5 := router{
		name:    "/users",
		pattern: "/users",
		handler: "/users handler here",
	}

	h6 := router{
		name:    "users-category",
		pattern: "/users/{category}",
		handler: "/{category} handler for the route /user[s]",
	}

	// r.router = hr
	r.node = map[rune]map[string]*path{
		117: {
			"user": &path{
				router: h1,
				node: map[rune]map[string]*path{
					47: {
						"{param}": &path{
							router: h2,
							node: map[rune]map[string]*path{
								105: {
									"images": &path{
										router: h3,
										node: map[rune]map[string]*path{
											47: {
												"{param}": &path{
													router: h4,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"users": &path{
				router: h5,
				node: map[rune]map[string]*path{
					47: {
						"{param}": &path{
							router: h6,
						},
					},
				},
			},
		},
	}

	rt, lr := r.getNode("/")

	a, la := r.getNode("/users/{category}")
	b, lb := r.getNode("/user/{name}")
	c, lc := r.getNode("/user/{name}/images")

	d, ld := r.getNode("/user/{name}/images/{id}")
	log.Printf("%v | lvl: %d \n", d.router.handler, ld)

	// if you change the handler, and thus it also changes in the tree
	// d.router.handler = "changed at the bottom"

	dc, ldc := r.getNode("/user/eugene/images/15")
	log.Printf("%v | lvl: %d \n", dc.router.handler, ldc)

	log.Printf("%v | lvl: %d \n", rt.router.handler, lr)
	log.Printf("%v | lvl: %d \n", a.router.handler, la)
	log.Printf("%v | lvl: %d \n", b.router.handler, lb)
	log.Printf("%v | lvl: %d \n", c.router.handler, lc)
	log.Printf("%v | lvl: %d \n", d.router.handler, ld)
}

func transverseNode(branch *path, part string) (*path, bool) {

	// get the charcode of the 1st character of the route part
	charcode := getCharCodeAt(part, 0)
	log.Printf(">> Character code: %v for word: %q\n", charcode, part)

	// if it a root route request ("")
	if charcode == 0 {
		log.Println("This is a root route request.")
		return branch, true
	}

	// get the pathmap of the Character code
	pathmap, ok := branch.node[charcode]

	// now check if it's node exists
	if ok {

		// if the node exists, then does it's path details exist
		if path, ok := pathmap[part]; ok {
			// if yes return it
			return path, true
		}

		// it's path details do not exist
		// return branch, false
	}

	log.Println("<========================= Route has not been found, maybe it is dymanic =========================>")

	// If no return has happened
	// then no map was found, try match for a {param} like pathmap
	// "{" charcode 47 as the holder
	if dynamicmap, ok := branch.node[47]; ok {
		log.Println("<========================= In here ===========================>")

		// if yes, the get the dynamic route pathmap map
		if path, ok := dynamicmap["{param}"]; ok {
			// if yes return it
			return path, true
		}
	}

	// If not return false, with the branch itself
	// false bool will aid in breaking the loop
	return branch, false
}

// get a specific route node
func (p *path) getNode(pattern string) (*path, int) {

	// If it is a route request
	if pattern == "/" {
		pattern = ""
	}

	// split the route pattern into parts "/" as the splitting key
	routeParts := strings.Split(pattern, "/")
	log.Printf("GET NODE: About to loop thru :-- %s | %d\n", routeParts, len(routeParts))

	var branchpath *path
	var pathFound bool
	var level int

	// Now loop thru the route pattern parts
	for lvl, part := range routeParts {

		// how deep are we/did we go
		// eg. /user is the 1st level after "/"
		level = lvl

		// the 1st loop will always be the root route path
		if part == "" && level == 0 {
			branchpath = p
		}

		// If a path is returned then it is the nested one
		// and fed back to get the next nested branch node
		// with the new/next in loop route pattern part
		// If not found break, thus return the last branch that was found
		// also the level tells us where the branch was found
		branchpath, pathFound = transverseNode(branchpath, part)
		log.Printf(">> Was the branch found: %v\n", pathFound)
		log.Printf(">> Branch returned: %q\n", branchpath)

		// If just one value exists in the array
		// It's the route array, return it
		if len(pattern) == 0 {
			log.Println(">> Only asked for route request.")
			return branchpath, level
		}

		// If no nested path is found, break it there is no match
		if !pathFound {
			// nothing found
			break
		}
	}

	// return nothing was found at all
	if branchpath == nil {
		return nil, -1
	}

	return branchpath, level
}

func (p *path) addNode(np *path) (*path, bool) {

	// Split the given route pattern into parts
	patternParts := strings.Split(np.pattern, "/")
	// patternParts := strings.Split(strings.Replace(np.pattern, "/", "", 1), "/")
	log.Printf("--------- Provided argument: %s, len: %v --------\n\n", patternParts, len(patternParts))

	// Prepare to piece them together one by one, adding to the end on each loop
	pieceTogether := make([]string, len(patternParts))

	var branchpath *path
	var isFound int

	// loop thru each route part checking if it exists
	// getting deeper into the tree after each single loop
	for level, part := range patternParts {

		// check if that node exists if it does, leave it alone, it not the one we want
		// loop into the next part
		pattern := strings.Join([]string{"/", part}, "")
		pieceTogether = append(pieceTogether, pattern)

		// join up the whole array into a route that can be used
		fullpath := strings.Replace(strings.Join(pieceTogether, ""), "//", "/", -1)
		log.Printf("------ Current path being checked :------ %v\n", fullpath)

		if part == "/" {
			branchpath = p
			continue
		}

		// now check if the node exists
		branchpath, isFound = p.getNode(fullpath)
		log.Printf("Branch found: %v | level: %v\n", branchpath, isFound)

		// If a branch was found
		if branchpath != nil && level < len(patternParts)-1 {
			log.Printf("--------- Branch exists!! --------\n", level)
			log.Printf("branchpath.node  len: %v\n", len(branchpath.node))

			// If the node of the branch does exist
			if branchpath.node == nil {
				log.Printf("Node does not exists!!\n", level)

				// The child part
				nextChildPart := patternParts[level+1]
				log.Printf("Next child path /%s\n", nextChildPart)

				// get the character code of the route part
				charcode := getCharCodeAt(nextChildPart, 0)

				// If not, creeate the node
				branchpath.node = make(map[rune]map[string]*path)

				// add a defualt route map
				routemap := make(map[string]*path)
				routemap[nextChildPart] = &path{
					router: router{
						name:    nextChildPart,
						pattern: fullpath,
						handler: strings.Join([]string{"/", nextChildPart, " handler for the route ", fullpath}, ""),
					},
					// node: make(map[rune]map[string]*path),
				}

				// NOTE: Do not reassign like This
				// branchpath = branchpath.node[charcode][nextChildPart]
				// you're overwriting it, and it will be detected on it own when looped
				branchpath.node[charcode] = routemap
			}

		}

		// which level are we in
		log.Printf("-------- Level: %v ---------\n\n", level)
	}

	return p, false
}

func (p *path) addRouteHandler(np *path) {
	p.router = np.router
}

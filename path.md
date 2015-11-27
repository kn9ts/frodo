## Custom Radix Tree/Search Algorithm Implementation for routing
__Used(or in implementation progress) in [Frodo, Go Web Framework repo](http://github.com/kn9ts/frodo)__

This is my hybrid of Radix trie/search algorithm. It is a modifcation of the known Radix trie algorithm eg. used in current(claimed) fastest Go routing package [httprouter](https://github.com/julienschmidt/httprouter)

It did facinate me but not enough. Radix tree algorithm unlike regular trees (where whole keys are compared en masse from their beginning up to the point of inequality), the key at each node is compared chunk-of-bits by chunk-of-bits, where the quantity of bits in that chunk at that node is the radix r of the radix trie. That's where I believe it fails. Comparing in chunck with for example a million incoming requests sounds resource intensive. But yet a far greater improvement from the old illiterating a single layer tree(an array of routes/paths) to find a match method - regular tree, beginner's luck I call it.

For example if we have the following route methods to add, see below:

__note:__ the _{param}_ signifies it's a dymanic route that can take in different names, id and e.t.c

```
/user/{name}
/user/profile
/user/settings
/user/{name}/image/{id}
/users
/users/scores
/images
/images/{id}
```

Since all we want to is to match are the words and params in between the slashes "/"

__Follows i the break down of my custom routing algorithm:__

- Split the route as "/" the splitting key, results be something like this:
```go
// let's the longest route: /user/{name}/image/{id}
// you end up with {"user", "{name}", "image", "{id}"} (Go Array)
routeParts := strings.Split(pattern, "/")
```
- Now loop through the route parts array, getting the character code of the 1st character for every route part, and matching for it in the route map eg.
```go
routePathMap, exists := RouteTree.node[charactercode] // exists will be false if it does not exist
```
- if it exists grab the node map and use that node to also verify the existence of the next route part using a recursive loop...going deeper for every route part
- The handlers and other meta data of the each route in question will also be stored side by side within the route part object map code extends the tree
- This give the space for existence for many routes nested in each other since they all share a common parent route part

So we end up hvaing something like this _(example in JS Object for ease in comprehension)_:

```js
var RouteTree = {
    "handler": {
        "pattern": "/",
        "handler": function() {},
        "isRegex": false,
        "depth": 0
    },
    "node": {
        105: {
            "user": {
                "handler": {
                    "pattern": "/",
                    "handler": function() {},
                    "isRegex": false,
                    "depth": 0
                },
                "node": {
                    97: {
                        "profile": {
                            "handler": {
                                "pattern": "/",
                                "handler": function() {},
                                "isRegex": false,
                                "depth": 0
                            },
                            "node": {}
                        },
                    },
                    100: {
                        "settings": {
                            "handler": {
                                "pattern": "/",
                                "handler": function() {},
                                "isRegex": false,
                                "depth": 0
                            },
                            "node": {}
                        }
                    },
                    45: {
                        "{name}": {
                            "handler": {
                                "pattern": "/",
                                "handler": function() {},
                                "isRegex": false,
                                "depth": 0
                            },
                            "node": {
                                85: {
                                    "image": {
                                        "handler": {
                                            "pattern": "/",
                                            "handler": function() {},
                                            "isRegex": false,
                                            "depth": 0
                                        },
                                        "node": {
                                            45: {
                                                "{id}": {
                                                    "handler": {
                                                        "pattern": "/",
                                                        "handler": function() {},
                                                        "isRegex": false,
                                                        "depth": 0,
                                                    },
                                                    "node": {}
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            },
            "users": {
                "handler": {
                    "pattern": "/",
                    "handler": function() {},
                    "isRegex": false,
                    "depth": 0
                },
                "node": {
                    100: {
                        "scores": {
                            "handler": {
                                "pattern": "/",
                                "handler": function() {},
                                "isRegex": false,
                                "depth": 0
                            },
                            "node": {}
                        }
                    }
                }
            }
        },
        85: {
            "images": {
                "handler": {
                    "pattern": "/",
                    "handler": function() {},
                    "isRegex": false,
                    "depth": 0
                },
                "node": {
                    45: {
                        "{id}": {
                            "handler": {
                                "pattern": "/",
                                "handler": function() {},
                                "isRegex": false,
                                "depth": 0
                            },
                            "node": {}
                        }
                    }
                }
            }
        }
    }
}

```

package main

import "fmt"

type Name struct {
	name string
}

func main() {
	var name *Name
	fmt.Println(name == nil)
}

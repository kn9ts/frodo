package main

import (
	"fmt"
	"github.com/kn9ts/frodo"
	// "net/http"
	"reflect"
)

/*

   ------- RESULTS: -------

   Value:  <*main.Parent Value>
   Type:  *main.Parent
   Kind:  ptr
   New:  <**main.Parent Value>
   PtrTo:  **main.Parent
   Interface:  &{{ }}
   Indirect:  <main.Parent Value>
   NumMethod:  10 false
   Elem(pointed struct, slice):  main.Parent
   Pointer ID:  833357997216

   "Default" "Delete" "Get" "Head" "MyMethod" "Options" "Patch" "Post" "Put" "_Construct"

*/

type Parent struct {
	Frodo.Controller
}

func (p *Parent) MyMethod() {

}

func main2() {

	ctrl := &Parent{}

	v := reflect.ValueOf(ctrl)
	fmt.Println("Value: ", v)
	fmt.Println("Type: ", v.Type())
	fmt.Println("Kind: ", v.Kind())
	fmt.Println("New: ", reflect.New(v.Type()))
	fmt.Println("PtrTo: ", reflect.PtrTo(v.Type()))
	fmt.Println("Interface: ", v.Interface())
	fmt.Println("Indirect: ", reflect.Indirect(v))
	fmt.Println("NumMethod: ", v.NumMethod(), v.CanAddr())

	if v.Kind().String() == "ptr" {
		// It panics if v's Kind is not Interface or Ptr.
		fmt.Println("Elem(pointed struct, slice): ", v.Type().Elem())

		// me := v.Interface()
		// fmt.Printf("%s\n", v.Type())

		// It panics if v's Kind is not Chan, Func, Map, Ptr, Slice, or UnsafePointer.
		fmt.Println("Pointer ID: ", v.Pointer())
	}

	for x := 0; x < v.NumMethod(); x++ {
		// (Func, Name, Index, PkgPath, Type)98
		fmt.Printf("%q\n", v.Type().Method(x).Name)
	}

	// for x := 0; x < v.NumField(); x++ {
	//  fmt.Println(v.Field(x))
	// }

}

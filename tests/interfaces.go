package main

import (
	"fmt"
	"github.com/kn9ts/frodo"
	"net/http"
	"reflect"
)

/*

	------- RESULTS: --------

	Value:  <*main.Name Value>
	Type:  *main.Name
	Kind:  ptr
	New:  <**main.Name Value>
	PtrTo:  **main.Name
	Interface:  &{{ }}
	Indirect:  <main.Name Value>
	NumMethod:  8 false
	Converted back, and posting this to alert you.
	Convertion Type: *main.Name
	&{{"" ""}}

	"Delete" "Get" "Head" "Options" "Patch" "Post" "Put" "_Construct"

*/

type Name struct {
	Frodo.Controller
}

type Profile struct {
	Frodo.ControllerInterface
}

func (n *Name) Index(w http.ResponseWriter, r *Frodo.Request) {
	fmt.Println("Converted back, and posting this to alert you.")
}

func (p *Profile) Get() {
	fmt.Println("Am method 'GET' for Profile.")
}

func testInterfaces() {
	fn(&Name{})
	// newf := func(w http.ResponseWriter, r *http.Request) {
	//
	// }
	// fn(newf)
}

func fn(args ...interface{}) {
	vl := reflect.ValueOf(args[0]).Type()
	makeHandler := func(h Frodo.HandleFunc) Frodo.HandleFunc {
		Frodo.Log.Debug("converting func(http.ResponseWriter, *Request) to Frodo.HandleFunc")
		return h
	}
	// Debug: First of check if it is a Frodo.HandleFunc type, might have been altered on first/previous loop
	// if not check the function if it suffices the HandleFunc type pattern
	// If it does -- func(http.ResponseWriter, *Request)
	// then convert it to a Frodo.HandleFunc type
	// this becomes neat since this what we expect to run
	// isHandleFunc := false
	if _, ok := args[0].(Frodo.HandleFunc); !ok {
		if value, ok := args[0].(func(http.ResponseWriter, *Frodo.Request)); ok && vl.Kind().String() == "func" {
			// morph it to it's dynamic data type
			args[0] = makeHandler(value)
			// isHandleFunc = true
		} else {
			// further checked if it is a Controller
			if _, isController := args[0].(Frodo.ControllerInterface); !isController {
				Frodo.Log.Fatal("Error: expected handler arguement provided to be an extension of Frodo.Controller or \"func(http.ResponseWriter, *Frodo.Request)\" type")
			}
			args[0] = args[0].(Frodo.ControllerInterface)
		}
	}
	me := args[0]

	v := reflect.ValueOf(me)
	fmt.Println("Value: ", v)
	fmt.Println("Type: ", v.Type())
	fmt.Println("Kind: ", v.Kind())
	fmt.Println("New: ", reflect.New(v.Type()))
	fmt.Println("PtrTo: ", reflect.PtrTo(v.Type()))
	fmt.Println("Interface: ", v.Interface())
	fmt.Println("Indirect: ", reflect.Indirect(v))
	fmt.Println("NumMethod: ", v.NumMethod(), v.CanAddr())

	abc := v.Interface().(Frodo.ControllerInterface)

	var w http.ResponseWriter
	var r = new(Frodo.Request)

	// This is if you knew which Method you are to invoke
	// in this case Name.Index
	// abc.Index(w, r)

	// If you dint, check for the method by it's name
	fn := v.MethodByName("Index")

	// Checking for a zero value
	if fn != (reflect.Value{}) {

		var newfn Frodo.HandleFunc

		// Then convert it back to HandleFunc
		// You have to know which type it is or are converting to
		if value, ok := fn.Interface().(func(http.ResponseWriter, *Frodo.Request)); ok && fn.Kind().String() == "func" {
			// morph it to it's dynamic data type
			newfn = makeHandler(value)
			// newfn(w, r)
			// isHandleFunc = true
		}

		// Then run it
		newfn(w, r)

		// Recheck
		fmt.Printf("What type is Index func: %q\n", reflect.ValueOf(newfn).Type().String())
	} else {
		// Recheck
		fmt.Printf("Index method Not found\n")
	}

	// panic: interface conversion: *main.Name is not ControllerInterface: missing method Post
	fmt.Printf("Convertion Type: %s\n", reflect.ValueOf(abc).Type())

	fmt.Printf("%q\n", abc)

	for x := 0; x < v.NumMethod(); x++ {
		// (Func, Name, Index, PkgPath, Type)98
		fmt.Printf("%q\n", v.Type().Method(x).Name)
	}

}

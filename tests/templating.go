package main

import (
	"fmt"
	"html/template"
	"os"
)

type Person struct {
	Name string
	Age  int
	Home string
}

func (p Person) Details() string {
	return fmt.Sprintf("Am %s, Am %d of age and live in %s", p.Name, p.Age, p.Home)
}

func templateExec() {
	Eugene := Person{"Eugene Mutai", 25, "Nairobi"}
	// --------------- or ----------------
	// Eugene := map[string]interface{}{
	// 	"Name": "Eugene Mutai",
	// 	"Age":  25,
	// 	"Home": "Nairobi",
	// }

	// tmpl, err := template.New("Foo").Parse("Am {{ .Name }}, Am {{ .Age }} of age and live in {{ .Home }}")
	tmpl, err := template.New("Foo").Parse("{{ .Details }}")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, Eugene)
	if err != nil {
		panic(err)
	}
}

func main() {
	templateExec()
}

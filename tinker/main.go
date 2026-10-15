package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type hand struct {
	name string
	age  int
}

func (h hand) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("").Parse(pageMarkup))
	tpl.Execute(w, fmt.Sprint(r.URL, " and ", h.age))
}

func newHandler(name string, options ...func(*hand)) hand {
	h := hand{
		name: name,
		age:  20, // default age
	}

	// each functional option modifies h and applies an option to it
	for _, option := range options {
		option(&h)
	}

	return h
}

func main() {
	fmt.Println("main")

	h := newHandler("John", handlerWithAdultAge)

	http.ListenAndServe(":8080", h)
}

var pageMarkup = `
<p>Hello World</p>
<p>{{.}}</p>
`

// functional options: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
// when passed to newHandler, it sets the age option
func handlerWithAdultAge(h *hand) {
	h.age = 50
}

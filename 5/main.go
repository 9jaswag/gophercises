package main

import (
	"fmt"
	hparse "hparse/parse"
	"strings"
)

func main() {
	sr := strings.NewReader(html)
	links, err := hparse.Parse(sr)
	if err != nil {
		fmt.Println("err")
	}

	fmt.Println(links)
}

/*
- only links on the same domain
- avoid cyclical links
- get domain via flag
*/
var html = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">
    A link to another page
    <span> with some nested span  </span>
  </a>
  <a href="/page-two">A link to a second page</a>
</body>
</html>
`

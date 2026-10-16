package main

import (
	"fmt"
	hparse "hparse/parse"
	"os"
)

func main() {
	f, err := os.Open("ex4.html")
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	links, err := hparse.Parse(f)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", links)
}

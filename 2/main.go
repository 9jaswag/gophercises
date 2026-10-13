package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	yamlFilePath := flag.String("yfile", "paths.yml", "file path to yaml file")
	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	urlPaths, err := parseYamlFile(*yamlFilePath)
	if err != nil {
		fmt.Println("Error reading yaml file", err)
	}

	yamlHandler, err := YAMLHandler(urlPaths, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func parseYamlFile(filePath string) ([]byte, error) {
	dataBSlice, err := os.ReadFile(filePath)
	if err != nil {
		return []byte{}, err
	}

	return dataBSlice, nil
}

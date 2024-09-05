package main

import (
	"flag"
	"fmt"
	"gophercises/urlshort"
	"io"
	"net/http"
	"os"
)

var YAMLpath string
var JSONpath string

func init() {
	flag.StringVar(&YAMLpath, "YAMLfile", "YAMLpath.yml", "a YAML file with url path")
	flag.StringVar(&JSONpath, "JSONfile", "JSONpath.json", "a JSON file with url path")
}

func main() {
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlFile, err := os.Open(YAMLpath)
	if err != nil {
		exit(fmt.Sprintf("Error opening file: %v\n", err))
	}
	defer yamlFile.Close()

	yaml, err := io.ReadAll(yamlFile)
	if err != nil {
		exit(fmt.Sprintf("Error reading file: %v\n", err))
	}

	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		fmt.Println("Error parsing YAML:", err)
		return
	}

	// Build the JSONHandler
	jsonFile, err := os.Open(JSONpath)
	if err != nil {
		exit(fmt.Sprintf("Error opening file: %v\n", err))
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		exit(fmt.Sprintf("Error reading file: %v\n", err))
	}

	jsonHandler, err := urlshort.JSONHandler([]byte(jsonData), mapHandler)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	_ = yamlHandler
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)

}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

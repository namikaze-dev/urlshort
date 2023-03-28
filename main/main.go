package main

import (
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/gophercises/urlshort"
)

type Options struct {
	YAML string
	JSON string
}

func main() {
	options := parseFlagOptions()
	YAMLData := readFile(options.YAML)
	JSONData := readFile(options.JSON)

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yamlHandler, err := urlshort.YAMLHandler(YAMLData, mapHandler)
	if err != nil {
		panic(err)
	}

	// Build the JSONHandler using the mapHandler as the
	// fallback
	jsonHandler, err := urlshort.JSONHandler(JSONData, yamlHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", jsonHandler)
}

func parseFlagOptions() Options {
	var options Options
	flag.StringVar(&options.YAML, "yaml", "", "yaml to load path/url combo from")
	flag.StringVar(&options.JSON, "json", "", "json to load path/url combo from")
	flag.Parse()
	return options
}

func readFile(fn string) []byte {
	data, err := fs.ReadFile(os.DirFS("."), fn)
	if err != nil {
		fmt.Printf("unexpected error while reading YAML %q: %v\n", fn, err)
		os.Exit(1)
	}
	return data
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v3"
)

type ParsedInput interface {
	Value() (path string, url string)
}

type ParsedYAML struct {
	Path string
	URL  string
}

func (py ParsedYAML) Value() (string, string) {
	return py.Path, py.URL
}

type ParsedJSON struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

func (pj ParsedJSON) Value() (string, string) {
	return pj.Path, pj.URL
}


// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if url, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, url, http.StatusSeeOther)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsedYAMLs []ParsedYAML
	err := yaml.Unmarshal(yml, &parsedYAMLs)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if url, ok := contains(parsedYAMLs, r.URL.Path); ok {
			http.Redirect(w, r, url, http.StatusSeeOther)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}, nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
// The only errors that can be returned all related to having
// invalid YAML data.
func JSONHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsedJSONs []ParsedJSON
	err := yaml.Unmarshal(yml, &parsedJSONs)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if url, ok := contains(parsedJSONs, r.URL.Path); ok {
			http.Redirect(w, r, url, http.StatusSeeOther)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}, nil
}

func contains[T ParsedInput](parsedInputs []T, path string) (string, bool) {
	for _, parsedInput := range parsedInputs {
		pPath, pURL := parsedInput.Value()
		if pPath == path {
			return pURL, true
		}
	}
	return "", false
}

package urlshortener

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

type UrlMapper struct {
	path string `yaml:"path"`
	url  string `yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if pathDest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, pathDest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
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
	urlMappers, err := parseByteYml(yml)
	if err != nil {
		return nil, err
	}
	pathUrlDictionary := buildDictionary(urlMappers)
	return MapHandler(pathUrlDictionary, fallback), nil
}

func buildDictionary(urlMappers []UrlMapper) map[string]string {
	pathUrlDictionary := make(map[string]string)
	for _, obj := range urlMappers {
		pathUrlDictionary[obj.path] = obj.url
	}
	return pathUrlDictionary
}

func parseByteYml(yml []byte) ([]UrlMapper, error) {
	var urlPathMapper []UrlMapper
	err := yaml.Unmarshal(yml, &urlPathMapper)
	if err != nil {
		return nil, err
	}
	return urlPathMapper, nil
}

package handlers

import (
	"io/ioutil"
	"net/http"
)

// Robots returns robots.txt.
func Robots() http.HandlerFunc {
	f, err := ioutil.ReadFile("robots.txt")
	if err != nil {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("User-agent: *\nDisallow: /"))
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(f)
	}
}

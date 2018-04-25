package handlers

import (
	"net/http"
)

// Favicon serves favicon.ico
func Favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/favicon.ico")
}

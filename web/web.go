package web

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"golang.org/x/text/language"
)

var langs = map[string]struct{}{
	"en": struct{}{},
	"fr": struct{}{},
}

var message = map[string]map[string]string{
	"en": {
		"greet": "It's rendering",
	},
	"fr": {
		"greet": "Affichage de",
	},
}

// Index handles index page.
func Index(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/html")

	// Set user preference language
	// by using cookies.
	if lang := r.URL.Query().Get("lang"); lang != "" {
		if _, ok := langs[lang]; ok {
			http.SetCookie(w, &http.Cookie{
				Name:     "lang",
				Value:    lang,
				Path:     "/",
				MaxAge:   604800,
				Secure:   true,
				HttpOnly: true,
			})
		}

		http.Redirect(w, r, r.URL.Path, 302)
		return
	}

	var lang string
	c, err := r.Cookie("lang")
	if err == nil {
		if _, ok := langs[c.Value]; ok {
			lang = c.Value
		} else {
			lang = "en"
		}
	} else {
		lg, _, _ := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
		if len(lg) < 1 {
			lang = "en"
		} else {
			switch lg[0].String() {
			case "fr-FR", "fr":
				lang = "fr"
			default:
				lang = "en"
			}
		}
	}

	t, err := template.New("index.html").Funcs(template.FuncMap{
		"T": func(s string) string {
			return message[lang][s]
		},
	}).ParseFiles("public/index.html", "public/head.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler("/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

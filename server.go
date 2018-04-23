package main

import (
	"net/http"
	"shlink-server/cmd"
	"shlink-server/handlers"
	middlewares "shlink-server/middlewares"
	"shlink-server/utils"
	"shlink-server/web"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

// initializeRoutes initialize api routes.
func initializeRoutes() *chi.Mux {
	// Router
	r := chi.NewRouter()

	// Enable CORS
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length",
			"Accept-Encoding", "X-CSRF-Token", "Authorization",
			"Accept-Language", "Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// Middlewares
	r.Use(middleware.StripSlashes)
	r.Use(cors.Handler)
	r.Use(middleware.RealIP)
	if cmd.Verbose {
		r.Use(middleware.Logger)
	}
	r.Use(middlewares.NewZapMiddleware("router", utils.Logger))
	r.Use(middleware.Recoverer)

	// Security header
	r.Use(middleware.SetHeader("X-XSS-Protection", "1; mode=block"))
	r.Use(middleware.SetHeader("X-Content-Type-Options", "nosniff"))
	r.Use(middleware.SetHeader("X-Frame-Options", "SAMEORIGIN"))

	// Endpoints
	r.Get("/", web.Index)
	r.Get("/robots.txt", handlers.Robots())
	r.Get("/{id}", handlers.RedirectURL)

	// API sub-router
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
		r.Use(middleware.NoCache)

		r.Get("/info/{id}", handlers.InfoURL)
		r.Get("/status", handlers.Status(version, platform, goVersion, goPlatform))
		r.Post("/create", handlers.Generate)
	})

	// Serve files
	web.FileServer(r, "/public", http.Dir("./public"))

	return r
}

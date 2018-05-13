package main

import (
	"fmt"
	"net/http"

	"github.com/afairon/shlink-server/handlers"
	middlewares "github.com/afairon/shlink-server/middlewares"
	"github.com/afairon/shlink-server/models"
	"github.com/afairon/shlink-server/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/cobra"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var cmdServer = &cobra.Command{
	Use:   "server",
	Short: "Start shlink http server",
	Run: func(cmd *cobra.Command, args []string) {
		// Setup zap logger
		logger := utils.SetupLogger(&lumberjack.Logger{
			Filename:   utils.Conf.Log.Filename,
			MaxSize:    utils.Conf.Log.MaxSize,
			MaxBackups: utils.Conf.Log.MaxBackups,
			MaxAge:     utils.Conf.Log.MaxAge,
			Compress:   true,
		})

		defer logger.Sync()

		// Connect to MongoDB
		session, err := models.Connect(utils.Conf)
		if err != nil {
			logger.Error(err.Error())
			panic(err)
		}

		defer session.Close()

		// Create indexes
		models.CreateIndexes()

		if !NoBanner {
			fmt.Printf("%s\n", Logo)
		}
		fmt.Printf("Shlink-Server %s\n", version)
		fmt.Printf("platform: %s\n", platform)
		fmt.Printf("go: %s\n", goVersion)
		fmt.Printf("built: %s\n", goPlatform)
		fmt.Printf("Listening on %s:%s\n", utils.Conf.Server.Host, utils.Conf.Server.Port)

		// Initialize routes
		r := initializeRoutes()

		// Start http server
		logger.Fatal(http.ListenAndServe(utils.Conf.Server.Host+":"+utils.Conf.Server.Port, r).Error())
	},
}

// initializeRoutes initialize api routes.
func initializeRoutes() (r *chi.Mux) {
	// Router
	r = chi.NewRouter()

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
	if Verbose {
		r.Use(middleware.Logger)
	}
	r.Use(middlewares.NewZapMiddleware("router", utils.Logger))
	r.Use(middleware.Recoverer)

	// Security header
	r.Use(middleware.SetHeader("X-XSS-Protection", "1; mode=block"))
	r.Use(middleware.SetHeader("X-Content-Type-Options", "nosniff"))
	r.Use(middleware.SetHeader("X-Frame-Options", "SAMEORIGIN"))

	// Endpoints
	r.Get("/{id}", handlers.RedirectURL)

	// API sub-router
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
		r.Use(middleware.NoCache)

		r.Get("/info/{id}", handlers.InfoURL)
		r.Get("/status", handlers.Status(version, platform, goVersion, goPlatform))
		r.Post("/create", handlers.Generate)
	})

	return
}

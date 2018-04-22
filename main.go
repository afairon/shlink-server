package main

import (
	"fmt"
	"net/http"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"

	"shlink-server/cmd"
	"shlink-server/models"
	utils "shlink-server/utils"
)

var (
	version    string
	platform   string
	goVersion  string
	goPlatform string
)

func main() {

	// Parsing commands and flags
	cmd.Execute()

	// Execute version
	if cmd.Version {
		if !cmd.NoBanner {
			fmt.Printf("%s\n", cmd.Logo)
		}
		fmt.Printf("Shlink-Server %s\n", version)
		fmt.Printf("platform: %s\n", platform)
		fmt.Printf("go: %s\n", goVersion)
		fmt.Printf("built: %s\n", goPlatform)

		os.Exit(0)
	}

	// Don't start server
	if !cmd.Start {
		os.Exit(0)
	}

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

	if !cmd.NoBanner {
		fmt.Printf("%s\n", cmd.Logo)
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
}

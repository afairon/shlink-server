package cmd

import (
	"fmt"
	"os"
	"shlink-server/utils"

	"github.com/spf13/cobra"
)

const (
	// Logo is the application banner
	Logo = `
  _____ _     _ _       _       _____                          
 / ____| |   | (_)     | |     / ____|                         
| (___ | |__ | |_ _ __ | | __ | (___   ___ _ ____   _____ _ __ 
 \___ \| '_ \| | | '_ \| |/ /  \___ \ / _ \ '__\ \ / / _ \ '__|
 ____) | | | | | | | | |   <   ____) |  __/ |   \ V /  __/ |   
|_____/|_| |_|_|_|_| |_|_|\_\ |_____/ \___|_|    \_/ \___|_|   
                                                                
`

	copyrights = `(C) Copyright 2018 Shlink. All Rights Reserved.`
)

var (
	// NoBanner will display banner.
	NoBanner bool

	// Start will start http server.
	Start bool

	// Version will print information about binary.
	Version bool

	// Verbose will enable go-chi verbose mode.
	Verbose bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "shlink-server",
	Short: Logo + copyrights,
}

var server = &cobra.Command{
	Use:   "server",
	Short: "Start shlink http server",
	Run: func(cmd *cobra.Command, args []string) {
		Start = true
	},
}

var version = &cobra.Command{
	Use:   "version",
	Short: "Shows binary version",
	Run: func(cmd *cobra.Command, args []string) {
		Version = true
	},
}

var config = &cobra.Command{
	Use:   "config",
	Short: "Shows config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[Database]")
		fmt.Printf("Host: %s\n", utils.Conf.Database.Host)
		fmt.Printf("Port: %s\n", utils.Conf.Database.Port)
		fmt.Printf("DB Name: %s\n", utils.Conf.Database.DB)
		fmt.Println("[Log]")
		fmt.Printf("Log Name: %s\n", utils.Conf.Log.Filename)
		fmt.Printf("Max Size: %d\n", utils.Conf.Log.MaxSize)
		fmt.Printf("Max Backups: %d\n", utils.Conf.Log.MaxBackups)
		fmt.Printf("Max Age: %d\n", utils.Conf.Log.MaxAge)
		fmt.Println("[Server]")
		fmt.Printf("Host: %s\n", utils.Conf.Server.Host)
		fmt.Printf("Port: %s\n", utils.Conf.Server.Port)
		fmt.Printf("Base: %s\n", utils.Conf.Server.Base)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	server.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose")

	RootCmd.PersistentFlags().BoolVar(&NoBanner, "no-banner", false, "Don't display banner")

	RootCmd.AddCommand(server)
	RootCmd.AddCommand(version)
	RootCmd.AddCommand(config)
}

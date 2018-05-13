package main

import (
	"fmt"
	"os"

	"github.com/afairon/shlink-server/utils"

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

	// Verbose will enable go-chi verbose mode.
	Verbose bool

	// version is the server version.
	version string

	// platform is the server platform.
	platform string

	// goVersion is the server compiled go version.
	goVersion string

	// goPlatform is the server compiled go platform.
	goPlatform string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "shlink",
	Short: Logo + copyrights,
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
	cmdServer.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose")

	RootCmd.PersistentFlags().BoolVar(&NoBanner, "no-banner", false, "Don't display banner")

	RootCmd.AddCommand(cmdServer)
	RootCmd.AddCommand(cmdVersion)
	RootCmd.AddCommand(config)
}

func main() {
	utils.Conf.ReadConfig()
	Execute()
}

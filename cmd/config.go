package main

import (
	"fmt"

	"github.com/afairon/shlink-server/utils"

	"github.com/spf13/cobra"
)

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

package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Shows binary version",
	Run: func(cmd *cobra.Command, args []string) {
		if !NoBanner {
			fmt.Printf("%s\n", Logo)
		}
		fmt.Printf("Shlink-Server %s\n", version)
		fmt.Printf("platform: %s\n", platform)
		fmt.Printf("go: %s\n", goVersion)
		fmt.Printf("built: %s\n", goPlatform)
	},
}

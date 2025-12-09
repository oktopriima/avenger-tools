package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "avenger-tools",
	Short: "Avenger Tools CLI",
	Long:  "A CLI tool to create and configure marvel new service boilerplates.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

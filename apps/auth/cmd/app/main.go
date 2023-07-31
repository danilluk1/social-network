package main

import (
	"os"

	"github.com/danilluk1/social-network/apps/auth/internal/app"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "Socialy auth microservice",
	Short:            "auth microservice",
	Long:             "This is socialy auth microservice",
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		app.NewApp().Run()
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

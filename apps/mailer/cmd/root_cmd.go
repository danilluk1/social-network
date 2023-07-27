package cmd

import (
	"context"
	"log"

	"github.com/danilluk1/social-network/apps/mailer/internal/conf"
	"github.com/spf13/cobra"
)

var configFile = ""

var rootCmd = cobra.Command{
	Long: "Root command of application",
	Run: func(cmd *cobra.Command, args []string) {
		serve(cmd.Context())
	},
}

func RootCommand() *cobra.Command {
	rootCmd.AddCommand(&versionCmd, &serveCmd)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "the config file to use")

	return &rootCmd
}

func loadConfig(ctx context.Context) *conf.Configuration {
	if ctx == nil {
		panic("context must not be nil")
	}

	config, err := conf.Load(configFile)
	if err != nil {
		log.Fatalf("failed to load configuration: %+v", err)
	}

	return config
}

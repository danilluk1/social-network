package cmd

import (
	"context"
	"log"

	"github.com/danilluk1/social-network/apps/auth/internal/conf"
	"github.com/spf13/cobra"
)

var configFile = ""

var rootCmd = cobra.Command{
	Use: "auth",
	Run: func(cmd *cobra.Command, args []string) {
		serve(cmd.Context())
	},
}

func RootCommand() *cobra.Command {
	rootCmd.AddCommand(&serveCmd, &versionCmd)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "the config file to use")

	return &rootCmd
}

func loadGlobalConfig(ctx context.Context) *conf.GlobalConfiguration {
	if ctx == nil {
		panic("context must not be nil")
	}

	config, err := conf.LoadGlobal(configFile)
	if err != nil {
		log.Fatalf("failed to load configuration: %+v", err)
	}

	return config
}

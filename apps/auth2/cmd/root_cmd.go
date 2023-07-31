package cmd

import (
	"context"
	"log"

	"github.com/danilluk1/social-network/apps/auth/internal/conf"
	"github.com/spf13/cobra"
)

var configFile = ""

var rootCmd = cobra.Command{
	Long: "Root point of application",
	Run: func(cmd *cobra.Command, args []string) {
		migrate(cmd, args)
		serve(cmd.Context())
	},
}

func RootCommand() *cobra.Command {
	rootCmd.AddCommand(&migrateCmd, &serveCmd, &versionCmd)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "the config file to use")

	return &rootCmd
}

func loadGlobalConfig(ctx context.Context) *conf.GlobalConfiguration {
	if ctx == nil {
		panic("context must not be nil")
	}

	config, err := conf.LoadGlobal(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %+v", err)
	}

	return config
}

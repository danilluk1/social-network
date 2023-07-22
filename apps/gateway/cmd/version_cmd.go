package cmd

import "github.com/spf13/cobra"

var serveCmd = cobra.Command{
  Use: "serve",
  Long: "Start Gateway",
  Run: func(cmd *cobra.Command, args []string)
}

func server(ctx context.Context) {
  
}

package cmd

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

var migrateCmd = cobra.Command{
	Long: "Run database migrations",
	Run:  migrate,
}

func migrate(cmd *cobra.Command, args []string) {
	globalConfig := loadGlobalConfig(cmd.Context())

	db, err := sql.Open("postgres", globalConfig.DB.URL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Set the dialect to "postgres"
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	// Run the migrations
	if err := goose.Up(db, globalConfig.DB.MigrationsPath); err != nil {
		panic(err)
	}

}

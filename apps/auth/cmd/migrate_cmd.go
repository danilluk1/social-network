package cmd

import (
	"database/sql"
	"embed"
	"net/url"
	"os"

	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var migrateCmd = cobra.Command{
	Use:  "auth",
	Long: "Migrate database structure. This will create new tables and add missing columns and indexes.",
	Run:  migrate,
}

func migrate(cmd *cobra.Command, args []string) {
	globalConfig := loadGlobalConfig(cmd.Context())

	if globalConfig.DB.Driver == "" && globalConfig.DB.URL != "" {
		u, err := url.Parse(globalConfig.DB.URL)
		if err != nil {
			logrus.Fatalf("%+v", errors.Wrap(err, "parsing db connection url"))
		}
		globalConfig.DB.Driver = u.Scheme
	}

	log := logrus.StandardLogger()

	if globalConfig.Logging.Level != "" {
		level, err := logrus.ParseLevel(globalConfig.Logging.Level)
		if err != nil {
			log.Fatalf("Failed to parse log level: %+v", err)
		}
		log.SetLevel(level)
		goose.SetLogger(log)

		if level != logrus.DebugLevel {
			goose.SetLogger(goose.NopLogger())
		}
	}

	var embedMigrations embed.FS
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		logrus.Fatalf("%+v", errors.Wrap(err, "opening db connection"))
	}
	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(globalConfig.DB.Driver); err != nil {
		logrus.Fatalf("%+v", errors.Wrap(err, "failed to set migration dialect"))
	}

	if err := goose.Up(db, globalConfig.DB.MigrationsPath); err != nil {
		logrus.Fatalf("%+v", errors.Wrap(err, "cannot run migrations"))
	}

	logrus.Infof("Auth migrations applied successfully")
}

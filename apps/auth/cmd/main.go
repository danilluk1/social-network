package main

import (
	"context"
	"os"

	"github.com/danilluk1/social-network/libs/config"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err)
	}

	if cfg.AppEnv == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := pgx.Connect(context.Background(), cfg.AuthPostgresUrl)
	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to database:")
	}
	defer conn.Close(context.Background())

	log.Info().Msg("Auth service running successfully🚀.")
}

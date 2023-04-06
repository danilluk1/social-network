package main

import (
	"os"

	"github.com/danilluk1/social-network/apps/gateway/internal/di"
	"github.com/danilluk1/social-network/apps/gateway/internal/middlewares"
	"github.com/danilluk1/social-network/libs/config"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize config")
	}

	if cfg.AppEnv == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	do.ProvideValue[config.Config](di.Provider, *cfg)

	errorMiddleware := middlewares.ErrorHandler()

	app := fiber.New(fiber.Config{
		ErrorHandler: errorMiddleware,
	})
	app.Use(cors.New())
	app.Use(compress.New())

	do.ProvideValue[auth.AuthClient](di.Provider, clients.NewAuth(cfg.AppEnv))

	// if cfg.AppEnv == "development" {
	// 	app.Get("/swagger/*", swagger.New(swagger.Config{}))
	// }

}

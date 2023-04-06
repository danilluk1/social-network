package main

import (
	"os"
	"reflect"
	"strings"

	"github.com/danilluk1/social-network/apps/gateway/internal/di"
	"github.com/danilluk1/social-network/apps/gateway/internal/middlewares"
	"github.com/danilluk1/social-network/apps/gateway/internal/services/redis"
	"github.com/danilluk1/social-network/apps/gateway/internal/types"
	"github.com/danilluk1/social-network/libs/config"
	"github.com/danilluk1/social-network/libs/grpc/clients"
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	rdb "github.com/redis/go-redis/v9"
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

	r := redis.New(cfg.RedisUrl)
	do.ProvideValue[*rdb.Client](di.Provider, r)

	storage := redis.NewCache(cfg.RedisUrl)

	validator := validator.New()
	en := en_US.New()
	uni := ut.New(en, en)
	transEN, _ := uni.GetTranslator("en_US")
	enTranslations.RegisterDefaultTranslations(validator, transEN)
	errorMiddleware := middlewares.ErrorHandler(transEN)
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	app := fiber.New(fiber.Config{
		ErrorHandler: errorMiddleware,
	})
	app.Use(cors.New())
	app.Use(compress.New())

	do.ProvideValue[auth.AuthClient](di.Provider, clients.NewAuth(cfg.AppEnv))

	v1 := app.Group("/v1")

	services := types.Services{
		RedisStorage:        storage,
		Validator:           validator,
		ValidatorTranslator: transEN,
	}

	// if cfg.AppEnv == "development" {
	// 	app.Get("/swagger/*", swagger.New(swagger.Config{}))
	// }

}

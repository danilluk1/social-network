package main

import (
	// "context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/danilluk1/social-network/apps/api/internal/graph"
	"github.com/danilluk1/social-network/apps/api/internal/graph/generated"
	"github.com/danilluk1/social-network/libs/config"
	"github.com/danilluk1/social-network/libs/grpc/clients"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"go.uber.org/zap"
)

func wrapHandler(f func(http.ResponseWriter, *http.Request)) func(ctx *fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		fasthttpadaptor.NewFastHTTPHandler(http.HandlerFunc(f))(ctx.Context())
	}
}

func main() {
	// ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	var logger *zap.Logger

	if cfg.AppEnv == "development" {
		l, _ := zap.NewDevelopment()
		logger = l
	} else {
		l, _ := zap.NewProduction()
		logger = l
	}

	zap.ReplaceGlobals(logger)

	app := fiber.New()

	authGrpc := clients.NewAuth(cfg.AppEnv)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		AuthGrpc: authGrpc,
	}}))

	// Serve GraphQL API
	app.Post("/graphql", func(c *fiber.Ctx) error {
		wrapHandler(h.ServeHTTP)(c)
		return nil
	})

	// Serve GraphQL Playground
	app.Get("/playground", func(c *fiber.Ctx) error {
		wrapHandler(playground.Handler("GraphQL", "/graphql"))(c)
		return nil
	})

	// Start the server
	err = app.Listen(":" + cfg.ApiGatewayPort)
	if err != nil {
		panic(err)
	}
	logger.Sugar().Info("Api started")
}

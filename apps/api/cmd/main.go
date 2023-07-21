package main

import (
	// "context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/danilluk1/social-network/apps/api/internal/graph"
	"github.com/danilluk1/social-network/apps/api/internal/graph/generated"
	"github.com/danilluk1/social-network/apps/api/internal/graph/middleware"
	"github.com/danilluk1/social-network/libs/config"
	"github.com/danilluk1/social-network/libs/grpc/clients"
	"github.com/go-chi/chi"
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

	router := chi.NewRouter()

	authGrpc := clients.NewAuth(cfg.AppEnv)

	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		AuthGrpc: authGrpc,
	}}))

	router.Use(middleware.Auth(authGrpc))

	// Serve GraphQL API
	router.Post("/graphql", h.ServeHTTP)

	// Serve GraphQL Playground
	router.Get("/playground", playground.Handler("GraphQL", "/graphql"))

	// Start the server
	logger.Sugar().Info("Api started")
	err = http.ListenAndServe(":"+cfg.ApiGatewayPort, router)
	if err != nil {
		panic(err)
	}
}

package graph

import (
	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	AuthGrpc auth.AuthClient
}

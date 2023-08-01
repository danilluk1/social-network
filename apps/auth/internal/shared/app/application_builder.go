package app

import (
	"github.com/danilluk1/social-network/libs/go/pkg/fxapp"
	"github.com/danilluk1/social-network/libs/go/pkg/fxapp/contracts"
)

type AuthApplicationBuilder struct {
	contracts.ApplicationBuilder
}

func NewAuthApplicationBuilder() *AuthApplicationBuilder {
	builder := &AuthApplicationBuilder{fxapp.NewApplicationBuilder()}

	return builder
}

func (a *AuthApplicationBuilder) Build() *AuthApplication {
	return NewAuthApplication(
		a.GetProvides(),
		a.GetDecorates(),
		a.Options(),
		a.Logger(),
		a.Environment(),
	)
}

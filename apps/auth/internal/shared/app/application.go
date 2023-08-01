package app

import (
	"github.com/danilluk1/social-network/apps/auth/internal/shared/configurations/auth"
	"github.com/danilluk1/social-network/libs/go/pkg/conf/environment"
	"github.com/danilluk1/social-network/libs/go/pkg/fxapp"
	"github.com/danilluk1/social-network/libs/go/pkg/logger"
	"go.uber.org/fx"
)

type AuthApplication struct {
	*auth.AuthServiceConfigurator
}

func NewAuthApplication(
	providers []interface{},
	decorates []interface{},
	options []fx.Option,
	logger logger.Logger,
	environment environment.Environment,
) *AuthApplication {
	app := fxapp.NewApplication(providers, decorates, options, logger, environment)
	return &AuthApplication{
		AuthServiceConfigurator: auth.NewAuthServiceConfigurator(app),
	}
}

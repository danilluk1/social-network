package contracts

import (
	"github.com/danilluk1/social-network/libs/go/pkg/conf/environment"
	"github.com/danilluk1/social-network/libs/go/pkg/logger"
	"go.uber.org/fx"
)

type ApplicationBuilder interface {
	ProvideModule(module fx.Option)
	Provide(constructors ...interface{})
	Decorate(constructors ...interface{})
	Build() Application

	GetProvides() []interface{}
	GetDecorates() []interface{}
	Options() []fx.Option
	Logger() logger.Logger
	Environment() environment.Environment
}

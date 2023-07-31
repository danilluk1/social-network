package conf

import (
	"github.com/danilluk1/social-network/libs/go/pkg/conf/environment"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"configfx",
	fx.Provide(func() environment.Environment {
		return environment.ConfigAppEnv()
	}),
)

var ModuleFunc = func(e environment.Environment) fx.Option {
	return fx.Module(
		"configfx",
		fx.Provide(func() environment.Environment {
			return e
		}),
	)
}

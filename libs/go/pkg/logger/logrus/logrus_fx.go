package logrus

import (
	"github.com/danilluk1/social-network/libs/go/pkg/logger"
	"github.com/danilluk1/social-network/libs/go/pkg/logger/conf"
	"go.uber.org/fx"
)

// Module provided to fxlog
// https://uber-go.github.io/fx/modules.html
var Module = fx.Module("logrusfx",
	// - order is not important in provide
	// - provide can have parameter and will resolve if registered
	// - execute its func only if it requested
	fx.Provide(
		fx.Annotate(
			NewLogrusLogger,
			fx.As(new(logger.Logger)),
		),
		conf.ProvideLogConfig,
	))

var ModuleFunc = func(l logger.Logger) fx.Option {
	return fx.Module("logrusfx",

		fx.Provide(conf.ProvideLogConfig),
		fx.Supply(fx.Annotate(l, fx.As(new(logger.Logger)))),
	)
}

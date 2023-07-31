package zap

import (
	"github.com/danilluk1/social-network/libs/go/pkg/logger"
	"github.com/danilluk1/social-network/libs/go/pkg/logger/conf"
	"go.uber.org/fx"
)

var Module = fx.Module("zapfx",

	fx.Provide(
		conf.ProvideLogConfig,
		NewZapLogger,
		fx.Annotate(
			NewZapLogger,
			fx.As(new(logger.Logger))),
	),
)

var ModuleFunc = func(l logger.Logger) fx.Option {
	return fx.Module(
		"zapfx",

		fx.Provide(conf.ProvideLogConfig),
		fx.Supply(fx.Annotate(l, fx.As(new(logger.Logger)))),
		fx.Supply(fx.Annotate(l, fx.As(new(ZapLogger)))),
	)
}

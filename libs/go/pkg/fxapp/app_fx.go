package fxapp

import (
	"time"

	"github.com/danilluk1/social-network/libs/go/pkg/conf"
	logConfig "github.com/danilluk1/social-network/libs/go/pkg/logger/conf"
	"github.com/danilluk1/social-network/libs/go/pkg/logger/external/fxlog"
	"github.com/danilluk1/social-network/libs/go/pkg/logger/logrus"
	"github.com/danilluk1/social-network/libs/go/pkg/logger/models"
	"github.com/danilluk1/social-network/libs/go/pkg/logger/zap"
	"go.uber.org/fx"
)

func CreateFxApp(
	app *application,
) *fx.App {
	var opts []fx.Option

	opts = append(opts, fx.Provide(app.provides...))

	opts = append(opts, fx.Decorate(app.decorates...))

	opts = append(opts, fx.Invoke(app.invokes...))

	app.options = append(app.options, opts...)

	AppModule := fx.Module("fxapp",
		app.options...,
	)

	var logModule fx.Option
	logOption, err := logConfig.ProvideLogConfig(app.environment)
	if err != nil || logOption == nil {
		logModule = zap.ModuleFunc(app.logger)
	} else if logOption.LogType == models.Logrus {
		logModule = logrus.ModuleFunc(app.logger)
	} else {
		logModule = zap.ModuleFunc(app.logger)
	}

	duration := 30 * time.Second

	fxApp := fx.New(
		fx.StartTimeout(duration),
		conf.ModuleFunc(app.environment),
		logModule,
		fxlog.FxLogger,
		fx.ErrorHook(NewFxErrorHandler(app.logger)),
		AppModule,
	)

	return fxApp
}

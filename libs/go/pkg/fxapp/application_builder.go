package fxapp

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/danilluk1/social-network/libs/go/pkg/conf/environment"
	"github.com/danilluk1/social-network/libs/go/pkg/constants"
	"github.com/danilluk1/social-network/libs/go/pkg/fxapp/contracts"
	"github.com/danilluk1/social-network/libs/go/pkg/logger"
	loggerConfig "github.com/danilluk1/social-network/libs/go/pkg/logger/conf"
	"github.com/danilluk1/social-network/libs/go/pkg/logger/models"
	"github.com/danilluk1/social-network/libs/go/pkg/logger/zap"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type applicationBuilder struct {
	provides    []interface{}
	decorates   []interface{}
	options     []fx.Option
	logger      logger.Logger
	environment environment.Environment
}

func NewApplicationBuilder(environments ...environment.Environment) contracts.ApplicationBuilder {
	env := environment.ConfigAppEnv(environments...)

	setConfigPath()
	var logger logger.Logger
	logoption, err := loggerConfig.ProvideLogConfig(env)
	if err != nil {
		log.Fatal(err)
	}
	if logoption == nil {
		logger = zap.NewZapLogger(logoption, env)
	} else if logoption.LogType == models.Logrus {
		//logger = logrous.NewLogrusLogger(logoption, env)
	} else {
		logger = zap.NewZapLogger(logoption, env)
	}

	return &applicationBuilder{logger: logger, environment: env}
}

func setConfigPath() {
	// https://stackoverflow.com/a/47785436/581476
	wd, _ := os.Getwd()

	// https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
	// get from `os env` or viper `internal registry`
	pn := viper.Get(constants.PROJECT_NAME_ENV)
	if pn == nil {
		return
	}
	for !strings.HasSuffix(wd, pn.(string)) {
		wd = filepath.Dir(wd)
	}
	// Get the absolute path of the executed project directory
	absCurrentDir, _ := filepath.Abs(wd)

	viper.Set(constants.AppRootPath, absCurrentDir)

	// Get the path to the "config" folder within the project directory
	configPath := filepath.Join(absCurrentDir, "config")
	viper.Set(constants.ConfigPath, configPath)
}

func (a *applicationBuilder) ProvideModule(module fx.Option) {
	a.options = append(a.options, module)
}

func (a *applicationBuilder) Provide(constructors ...interface{}) {
	a.provides = append(a.provides, constructors...)
}

func (a *applicationBuilder) Decorate(constructors ...interface{}) {
	a.decorates = append(a.decorates, constructors...)
}

func (a *applicationBuilder) Build() contracts.Application {
	app := NewApplication(a.provides, a.decorates, a.options, a.logger, a.environment)

	return app
}

func (a *applicationBuilder) GetProvides() []interface{} {
	return a.provides
}

func (a *applicationBuilder) GetDecorates() []interface{} {
	return a.decorates
}

func (a *applicationBuilder) Options() []fx.Option {
	return a.options
}

func (a *applicationBuilder) Logger() logger.Logger {
	return a.logger
}

func (a *applicationBuilder) Environment() environment.Environment {
	return a.environment
}

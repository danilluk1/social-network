package contracts

import (
	"context"

	"github.com/danilluk1/social-network/libs/go/pkg/conf/environment"
	"github.com/danilluk1/social-network/libs/go/pkg/logger"
	"go.uber.org/fx"
)

type Application interface {
	Container
	RegisterHook(function interface{})
	Run()
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Wait() <-chan fx.ShutdownSignal
	Logger() logger.Logger
	Environment() environment.Environment
}

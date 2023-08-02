package auth

import (
	"github.com/danilluk1/social-network/apps/auth/internal/shared/configurations/auth/infrastructure"
	"github.com/danilluk1/social-network/libs/go/pkg/fxapp/contracts"
)

type AuthServiceConfigurator struct {
	contracts.Application
	infrastructureConfigurator *infrastructure.InfrastructureConfigurator
}

func NewAuthServiceConfigurator(app contracts.Application) *AuthServiceConfigurator {
	infraConfigurator := infrastructure.NewInfrastructureConfigurator(app)

	return &AuthServiceConfigurator{
		Application:                app,
		infrastructureConfigurator: infraConfigurator,
	}
}

func (ic *AuthServiceConfigurator) Configure() {
	ic.infrastructureConfigurator.ConfigInfrastructures()

	// ic.ResolveFunc(func() error {
	//   err := ic.
	// })
}

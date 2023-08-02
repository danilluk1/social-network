package app

import "github.com/danilluk1/social-network/apps/auth/internal/shared/configurations/auth"

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	appBuilder := NewAuthApplicationBuilder()
	appBuilder.ProvideModule(auth.AuthServiceModule)

	app := appBuilder.Build()
	app.Configure()

	app.Logger().Info("Starting auth microservice")
	app.Run()
}

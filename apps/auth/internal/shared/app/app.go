package app

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() {
	appBuilder := NewAuthApplicationBuilder()
	appBuilder.ProvideModule()
}

package app

import (
	"creditcard/internal/handler/cli"
)

type App struct {
	handler cli.IHandler
}

func NewApp() *App {
	return &App{
		handler: cli.NewHandler(),
	}
}

// Run TODO: args may be more than 1, so maybe we need to firstly get command, then check len of arguments
func (a *App) Run() error {
	return a.handler.Gateway()
}

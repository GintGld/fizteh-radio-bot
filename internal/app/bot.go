package app

import "log/slog"

type App struct {
	log     *slog.Logger
	tgToken string
	yaToken string
}

func New(
	log *slog.Logger,
	tgToken string,
	yaToken string,
) *App {
	return &App{
		log:     log,
		tgToken: tgToken,
		yaToken: yaToken,
	}
}

func (a *App) Run() error {
	panic("not implemented")
}

func (a *App) Stop() error {
	panic("mot implemented")
}

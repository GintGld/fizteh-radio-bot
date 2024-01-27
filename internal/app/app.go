package app

import (
	"log/slog"

	tele "gopkg.in/telebot.v3"

	"github.com/GintGld/fizteh-radio-bot/internal/lib/logger/sl"
)

type App struct {
	log *slog.Logger
	bot *tele.Bot
}

// New returns new App instance
func New(
	log *slog.Logger,
	token string,
) *App {
	bot, err := tele.NewBot(tele.Settings{
		Token: token,
	})
	if err != nil {
		log.Error("failed to create bot", sl.Err(err))
		panic("failed to create bot " + err.Error())
	}

	// TODO: mount all stuff here

	return &App{
		log: log,
		bot: bot,
	}
}

// Run starts bot serving
func (a *App) Run() {
	a.bot.Start()
}

func (a *App) Stop() {
	a.bot.Stop()
}

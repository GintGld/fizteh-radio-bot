package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-telegram/bot"
)

type App struct {
	log     *slog.Logger
	bot     *bot.Bot
	yaToken string

	server *http.Server
	cancel context.CancelFunc
}

// New returns new bot instance.
func New(
	log *slog.Logger,
	tgToken string,
	yaToken string,
	webhookAddr string,
) *App {
	bot, err := bot.New(tgToken)
	if err != nil {
		panic("failed to create bot")
	}

	server := &http.Server{
		Addr:    webhookAddr,
		Handler: bot.WebhookHandler(),
	}

	return &App{
		log:     log,
		bot:     bot,
		server:  server,
		yaToken: yaToken,
	}
}

// Run starts bot with webhook.
func (a *App) Run(ctx context.Context) error {
	ctx, a.cancel = context.WithCancel(ctx)

	go a.bot.StartWebhook(ctx)

	return a.server.ListenAndServe()
}

// Stop stops bot and its wekhook server.
func (a *App) Stop() error {
	a.cancel()

	return a.server.Close()
}

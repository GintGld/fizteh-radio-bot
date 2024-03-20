package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/GintGld/fizteh-radio-bot/internal/app"
	"github.com/GintGld/fizteh-radio-bot/internal/config"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/logger/slogpretty"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env, cfg.LogPath)

	log.Info("start bot", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Create bot instance
	app := app.New(
		log,
		getTelegramToken(),
		cfg.RadioAddr,
		getYandexToken(),
		cfg.WebhookAddr,
		cfg.TmpDir,
		cfg.UseFiller,
	)

	// Run bot
	go app.Run(context.Background())

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.Stop()

	log.Info("Gracefully stopped")
}

func setupLogger(env, logPath string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envProd:
		var logWriter io.Writer

		if logPath == "" {
			logWriter = os.Stdout
		} else {
			var err error
			logWriter, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				panic("failed to open log file. Error: " + err.Error())
			}
		}

		log = slog.New(
			slog.NewJSONHandler(logWriter, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		panic("unknown environment " + env)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

func getTelegramToken() string {
	token := os.Getenv("TG_TOKEN")

	if token == "" {
		panic("telegram token not specified")
	}

	return token
}

func getYandexToken() string {
	token := os.Getenv("YA_TOKEN")

	if token == "" {
		panic("telegram token not specified")
	}

	return token
}

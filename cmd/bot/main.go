package main

import (
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

	log.Info("starting bot", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// Start bot
	app := app.New(
		log,
		getToken(),
	)

	go app.Run()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.Stop()
	log.Info("gracefully stopped")
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

func getToken() string {
	token := os.Getenv("TOKEN")

	if token == "" {
		panic("token not specified")
	}

	return token
}

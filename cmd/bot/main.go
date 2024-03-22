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

func main() {
	cfg := config.MustLoad()

	logSrv := setupLogger(cfg.Log.Srv.Pretty, cfg.Log.Srv.Level, cfg.Log.Srv.Path)
	logTg := setupLogger(cfg.Log.Tg.Pretty, cfg.Log.Tg.Level, cfg.Log.Tg.Path)

	logSrv.Info("start bot", slog.String("env", cfg.Env))
	logSrv.Debug("debug messages are enabled")

	// Create bot instance
	app := app.New(
		logSrv,
		logTg,
		getTelegramToken(),
		cfg.RadioAddr,
		getYandexToken(),
		cfg.WebhookAddr,
		cfg.TmpDir,
		cfg.UserCacheFile,
		cfg.UseFiller,
	)

	// Run bot
	go app.Run(context.Background())

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.Stop()

	logSrv.Info("Gracefully stopped")
}

func setupLogger(pretty bool, level slog.Level, logPath string) *slog.Logger {
	var log *slog.Logger

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

	switch pretty {
	case true:
		log = setupPrettySlog(logWriter, level)
	case false:
		log = slog.New(
			slog.NewJSONHandler(logWriter, &slog.HandlerOptions{Level: level}),
		)
	}

	return log
}

func setupPrettySlog(writer io.Writer, level slog.Level) *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: level,
		},
	}

	handler := opts.NewPrettyHandler(writer)

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

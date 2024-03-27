package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/GintGld/fizteh-radio-bot/internal/controller/autodj"
	"github.com/GintGld/fizteh-radio-bot/internal/controller/datetime"
	"github.com/GintGld/fizteh-radio-bot/internal/controller/help"
	"github.com/GintGld/fizteh-radio-bot/internal/controller/schedule"
	"github.com/GintGld/fizteh-radio-bot/internal/controller/search"
	"github.com/GintGld/fizteh-radio-bot/internal/controller/start"
	"github.com/GintGld/fizteh-radio-bot/internal/controller/upload"

	authSrv "github.com/GintGld/fizteh-radio-bot/internal/service/auth"
	"github.com/GintGld/fizteh-radio-bot/internal/service/filler"
	libSrv "github.com/GintGld/fizteh-radio-bot/internal/service/library"
	schSrv "github.com/GintGld/fizteh-radio-bot/internal/service/schedule"
	"github.com/GintGld/fizteh-radio-bot/internal/service/session"

	radioCl "github.com/GintGld/fizteh-radio-bot/internal/client/radio"
	yandexCl "github.com/GintGld/fizteh-radio-bot/internal/client/yandex"
)

type App struct {
	log *slog.Logger
	bot *bot.Bot

	server *http.Server
	cancel context.CancelFunc
}

// New returns new bot instance.
func New(
	logSrv *slog.Logger,
	logTg *slog.Logger,
	tgToken string,
	radioAddr string,
	yaToken string,
	webhookAddr string,
	tmpDir string,
	userCacheFile string,
	srvFiller bool,
) *App {
	// default handlers
	errorHandler := getErrorHandler(logTg)
	defaultHandler := getDefaultHandler(logTg, errorHandler)

	bot, err := bot.New(tgToken,
		bot.WithDefaultHandler(defaultHandler),
	)
	if err != nil {
		panic("failed to create bot: " + err.Error())
	}

	// Clients
	var (
		authClient        authSrv.AuthClient
		libClient         libSrv.LibraryClient
		libGetMediaClient schSrv.LibraryClient
		yaClient          libSrv.YaClient
		schClient         schSrv.ScheduleClient
		djClient          schSrv.AutoDJClient
	)

	radioClient := radioCl.New(
		radioAddr,
	)
	yandexClient := yandexCl.New(
		yaToken,
		tmpDir,
	)

	authClient = radioClient
	libClient = radioClient
	libGetMediaClient = radioClient
	yaClient = yandexClient
	schClient = radioClient
	djClient = radioClient

	// Services
	var (
		auth           start.Auth
		libSearchSrv   search.Library
		scheduleAddSrv datetime.ScheduleAdd
		mediaUploadSrv upload.MediaUpload
		getScheduleSrv schedule.Schedule
		dj             autodj.AutoDJ
	)

	if srvFiller {
		filler := filler.New()

		auth = filler
		libSearchSrv = filler
		scheduleAddSrv = filler
		mediaUploadSrv = filler
		getScheduleSrv = filler
		dj = filler
	} else {
		a := authSrv.New(
			logSrv,
			authClient,
			userCacheFile,
		)
		l := libSrv.New(
			logSrv,
			a,
			libClient,
			yaClient,
		)
		s := schSrv.New(
			logSrv,
			a,
			libGetMediaClient,
			schClient,
			djClient,
		)

		auth = a
		libSearchSrv = l
		scheduleAddSrv = s
		mediaUploadSrv = l
		getScheduleSrv = s
		dj = s
	}

	// routing
	session := session.New[string]()

	router := ctr.NewRouter(
		bot, session,
	)

	start.Register(
		router.With("start"),
		auth,
		session,
		errorHandler,
	)
	help.Register(
		router.With("help"),
		auth,
		errorHandler,
	)
	search.Register(
		router.With("lib"),
		auth,
		libSearchSrv,
		scheduleAddSrv,
		session,
		errorHandler,
	)
	upload.Register(
		router.With("upload"),
		auth,
		mediaUploadSrv,
		session,
		errorHandler,
		tmpDir,
	)
	schedule.Register(
		router.With("sch"),
		auth,
		getScheduleSrv,
		session,
		errorHandler,
	)
	autodj.Register(
		router.With("dj"),
		auth,
		dj,
		session,
		errorHandler,
	)

	return &App{
		log: logSrv,
		bot: bot,
		server: &http.Server{
			Addr:    webhookAddr,
			Handler: bot.WebhookHandler(),
		},
	}
}

func getDefaultHandler(log *slog.Logger, errorHandler bot.ErrorsHandler) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "defaultHandler"

	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message != nil {
			if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   ctr.UnexpectedMsg,
			}); err != nil {
				chatId := update.Message.Chat.ID
				errorHandler(fmt.Errorf("%s [%d]: %w", op, chatId, err))
				return
			}
			return
		}

		if update.CallbackQuery != nil {
			chatId := update.CallbackQuery.Message.Message.Chat.ID
			ok, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
			})
			if err != nil {
				errorHandler(fmt.Errorf("%s [%d]: %w", op, chatId, err))
				return
			}
			if !ok {
				errorHandler(fmt.Errorf("%s [%d]: %s", op, chatId, "callback answer failed"))
			}

			if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.CallbackQuery.From.ID,
				Text:   ctr.UndefMsg,
			}); err != nil {
				errorHandler(fmt.Errorf("%s [%d]: %w", op, chatId, err))
			}

			log.Error(
				"unexpected callback",
				slog.Int("id", int(chatId)),
				slog.String("callback", update.CallbackQuery.Data),
			)

			return
		}

	}
}

func getErrorHandler(log *slog.Logger) bot.ErrorsHandler {
	return func(err error) {
		log.Error(err.Error())
	}
}

// Run starts bot with webhook.
func (a *App) Run(ctx context.Context) error {
	// FIXME webhook
	// TODO add to config wekhook and update options
	ctx, a.cancel = context.WithCancel(ctx)
	// go a.bot.StartWebhook(ctx)
	// return a.server.ListenAndServe()

	go a.bot.Start(ctx)

	return nil
}

// Stop stops bot and its wekhook server.
func (a *App) Stop() error {
	a.cancel()
	return nil
	// return a.server.Close()
}

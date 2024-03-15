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
	tmpDir string,
	useFiller bool,
) *App {
	bot, err := bot.New(tgToken,
		bot.WithDefaultHandler(defaultHandler),
	)
	if err != nil {
		panic("failed to create bot: " + err.Error())
	}

	var (
		auth           start.Auth
		libSearchSrv   search.LibrarySearch
		scheduleAddSrv datetime.ScheduleAdd
		mediaUploadSrv upload.MediaUpload
		getScheduleSrv schedule.Schedule
		dj             autodj.AutoDJ
	)

	if useFiller {
		filler := filler.New()

		auth = filler
		libSearchSrv = filler
		scheduleAddSrv = filler
		mediaUploadSrv = filler
		getScheduleSrv = filler
		dj = filler
	} else {
		a := authSrv.New(
			log,
			nil, // TODO
		)
		l := libSrv.New(
			log,
			nil, // TODO
			nil, // TODO
			nil, // TODO
		)
		s := schSrv.New(
			log,
			nil, // TODO
			nil, // TODO
		)

		auth = a
		libSearchSrv = l
		scheduleAddSrv = s
		mediaUploadSrv = l
		getScheduleSrv = s
		dj = s

		panic("not implemented") // FIXME add client to services
	}

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
		session,
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
		log: log,
		bot: bot,
		server: &http.Server{
			Addr:    webhookAddr,
			Handler: bot.WebhookHandler(),
		},
		yaToken: yaToken,
	}
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message != nil {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.From.ID,
			Text:   ctr.UndefMsg,
		}); err != nil {
			errorHandler(err)
		}
	}

	if update.CallbackQuery != nil {
		ok, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})
		if err != nil {
			errorHandler(err)
			return
		}
		if !ok {
			errorHandler(fmt.Errorf("callback answer failed"))
		}

		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.From.ID,
			Text:   ctr.UndefMsg,
		}); err != nil {
			errorHandler(err)
		}
	}
}

// TODO add another slog.Logger
// dump it to special file.
func errorHandler(err error) {
	fmt.Println(err.Error())
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

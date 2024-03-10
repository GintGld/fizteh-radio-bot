package library

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/GintGld/fizteh-radio-bot/internal/controller/search"
	"github.com/GintGld/fizteh-radio-bot/internal/controller/upload"
)

const (
	cmdSearch = "search"
	cmdUpload = "upload"
)

type library struct {
	router  *ctr.Router
	auth    Auth
	session ctr.Session
	onError bot.ErrorsHandler
}

type Auth interface {
	IsKnown(id int64) bool
	Login(login, pass string) error
}

func Register(
	router *ctr.Router,
	auth Auth,
	libSearch search.LibrarySearch,
	scheduleAdd search.ScheduleAdd,
	mediaUpload upload.MediaUpload,
	session ctr.Session,
	onError bot.ErrorsHandler,
	tmpDir string,
) {
	l := library{
		router:  router,
		auth:    auth,
		session: session,
		onError: onError,
	}

	router.RegisterCommand(l.libraryMainMenu)

	search.Register(
		router.With("search"),
		libSearch,
		scheduleAdd,
		session,
		l.cancelSubmodule,
		onError,
	)

	upload.Register(
		router.With("upload"),
		mediaUpload,
		session,
		l.cancelSubmodule,
		onError,
		tmpDir,
	)
}

func (l *library) libraryMainMenu(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID

	if !l.auth.IsKnown(userId) {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrUnknown,
		}); err != nil {
			l.onError(err)
		}
		return
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        ctr.LibMainMenuMessage,
		ReplyMarkup: l.mainMenuMarkup(),
	}); err != nil {
		l.onError(err)
	}
}

func (l *library) cancelSubmodule(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	chatId := mes.Message.Chat.ID

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        ctr.LibMainMenuMessage,
		ReplyMarkup: l.mainMenuMarkup(),
	}); err != nil {
		l.onError(err)
	}
}

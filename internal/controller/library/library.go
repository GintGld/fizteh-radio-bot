package library

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
)

const (
	base = "/lib"

	cmdSearch = base + "/search"
	cmdUpload = base + "/upload"
)

type library struct {
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
	session ctr.Session,
	onError bot.ErrorsHandler,
) {
	l := library{
		auth:    auth,
		session: session,
		onError: onError,
	}

	router.RegisterCommand(base, l.libraryMainMenu)
	router.RegisterCallback(cmdSearch, l.search)
	router.RegisterCallback(cmdUpload, l.upload)
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

func (l *library) search(ctx context.Context, b *bot.Bot, update *models.Update) {
	// TODO
}

func (l *library) upload(ctx context.Context, b *bot.Bot, update *models.Update) {
	// TODO
}

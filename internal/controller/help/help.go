package help

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
)

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
	router.RegisterCommand("/help", func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.Message.From.ID
		chatId := update.Message.Chat.ID

		if !auth.IsKnown(userId) {
			if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatId,
				Text:   ctr.ErrUnknown,
			}); err != nil {
				onError(err)
			}
			return
		}

		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.HelpMessage,
		}); err != nil {
			onError(err)
		}
	})
}

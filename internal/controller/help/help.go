package help

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
)

type Auth interface {
	IsKnown(ctx context.Context, id int64) bool
}

func Register(
	router *ctr.Router,
	auth Auth,
	session ctr.Session,
	onError bot.ErrorsHandler,
) {
	router.RegisterCommand(func(ctx context.Context, b *bot.Bot, update *models.Update) {

		chatId := update.Message.Chat.ID

		if !auth.IsKnown(ctx, chatId) {
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

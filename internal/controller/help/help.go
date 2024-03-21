package help

import (
	"context"
	"fmt"

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
	onError bot.ErrorsHandler,
) {
	router.RegisterCommand(func(ctx context.Context, b *bot.Bot, update *models.Update) {
		const op = "help"

		chatId := update.Message.Chat.ID

		if !auth.IsKnown(ctx, chatId) {
			if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatId,
				Text:   ctr.ErrUnknown,
			}); err != nil {
				onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
			}
			return
		}

		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.HelpMessage,
		}); err != nil {
			onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
	})
}

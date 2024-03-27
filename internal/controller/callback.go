package controller

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type CallbackAnswerer struct {
	onError bot.ErrorsHandler
}

func (c *CallbackAnswerer) CallbackAnswer(ctx context.Context, b *bot.Bot, callbackQuery *models.CallbackQuery) {
	const op = "CallbackAnswerer.callbackAnswer"

	chatId := callbackQuery.Message.Message.Chat.ID

	ok, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
	})
	if err != nil {
		c.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		return
	}
	if !ok {
		c.onError(fmt.Errorf("%s [%d]: %s", op, chatId, "callback answer failed"))
	}
}

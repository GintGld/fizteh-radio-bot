package search

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (s *Search) nextMediaSlide(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	id := s.mediaPage.Get(userId) + 1
	s.mediaPage.Set(userId, id)

	res := s.mediaResults.Get(userId)
	s.mediaSelected.Set(userId, res[id-1])

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        s.mediaRepr(res[id-1]),
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: s.mediaSliderMarkup(id, len(res)),
	}); err != nil {
		s.onError(err)
	}
}

func (s *Search) prevMediaSlide(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	id := s.mediaPage.Get(userId) - 1
	s.mediaPage.Set(userId, id)

	res := s.mediaResults.Get(userId)
	s.mediaSelected.Set(userId, res[id-1])

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        s.mediaRepr(res[id-1]),
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: s.mediaSliderMarkup(id, len(res)),
	}); err != nil {
		s.onError(err)
	}
}

func (s *Search) canceledDateTimeSelector(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	userId := mes.Message.From.ID
	chatId := mes.Message.Chat.ID

	id := s.mediaPage.Get(userId)
	res := s.mediaResults.Get(userId)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        s.mediaRepr(res[id-1]),
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: s.mediaSliderMarkup(id, len(res)),
	}); err != nil {
		s.onError(err)
	}
}

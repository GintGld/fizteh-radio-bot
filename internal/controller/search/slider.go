package search

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
)

func (s *Search) updateSlide(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID
	direction := s.router.GetState(update.CallbackQuery.Data)

	id := s.mediaPage.Get(chatId)
	switch direction {
	case "prev":
		id--
	case "next":
		id++
	default:
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			s.onError(err)
		}
	}
	s.mediaPage.Set(chatId, id)

	res := s.mediaResults.Get(chatId)
	s.mediaSelected.Set(chatId, res[id-1])

	msgId := s.msgIdStorage.Get(chatId)

	msg, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   msgId,
		Text:        s.mediaRepr(res[id-1]),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: s.mediaSliderMarkup(id, len(res)),
	})
	if err != nil {
		s.onError(err)
	}

	s.msgIdStorage.Set(chatId, msg.ID)
}

func (s *Search) canceledDateTimeSelector(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	chatId := mes.Message.Chat.ID

	id := s.mediaPage.Get(chatId)
	res := s.mediaResults.Get(chatId)

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        s.mediaRepr(res[id-1]),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: s.mediaSliderMarkup(id, len(res)),
	})
	if err != nil {
		s.onError(err)
	}
	s.msgIdStorage.Set(chatId, msg.ID)
}

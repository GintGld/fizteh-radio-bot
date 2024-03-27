package search

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
)

func (s *search) updateSlide(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "search.updateSlide"

	s.CallbackAnswer(ctx, b, update.CallbackQuery)

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
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
	}
	s.mediaPage.Set(chatId, id)

	res := s.mediaResults.Get(chatId)
	s.mediaSelectedStorage.Set(chatId, res[id-1])

	msgId := s.msgIdStorage.Get(chatId)

	msg, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   msgId,
		Text:        res[id-1].String(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: s.mediaSliderMarkup(id, len(res)),
	})
	if err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}

	s.msgIdStorage.Set(chatId, msg.ID)
}

func (s *search) canceledDateTimeSelector(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	const op = "search.canceledDateTimeSelector"

	chatId := mes.Message.Chat.ID

	id := s.mediaPage.Get(chatId)
	res := s.mediaResults.Get(chatId)

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        res[id-1].String(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: s.mediaSliderMarkup(id, len(res)),
	})
	if err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
	s.msgIdStorage.Set(chatId, msg.ID)
}

func (s *search) updateMedia(ctx context.Context, b *bot.Bot, msg models.MaybeInaccessibleMessage) {
	const op = "search.updateMedia"

	chatId := msg.Message.Chat.ID

	conf := s.mediaSelectedStorage.Get(chatId)

	if err := s.lib.UpdateMedia(ctx, chatId, conf); err != nil {
		// TODO handle errors
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    chatId,
		MessageID: s.msgIdStorage.Get(chatId),
		Text:      ctr.LibSearchUpdatedSuccess,
	}); err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (s *search) closedUpdateMedia(ctx context.Context, b *bot.Bot, msg models.MaybeInaccessibleMessage) {
	const op = "search.closedUpdateMedia"

	chatId := msg.Message.Chat.ID

	id := s.mediaPage.Get(chatId)
	res := s.mediaResults.Get(chatId)

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   s.msgIdStorage.Get(chatId),
		Text:        res[id-1].String(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: s.mediaSliderMarkup(id, len(res)),
	}); err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

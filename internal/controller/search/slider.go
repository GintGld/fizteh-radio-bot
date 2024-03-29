package search

import (
	"context"
	"fmt"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
)

func (s *search) updateSlide(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "search.updateSlide"

	s.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID
	direction := s.router.GetState(update.CallbackQuery.Data)

	id := s.mediaPageStorage.Get(chatId)
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
	s.mediaPageStorage.Set(chatId, id)

	res := s.mediaResultsStorage.Get(chatId)
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

	id := s.mediaPageStorage.Get(chatId)
	res := s.mediaResultsStorage.Get(chatId)

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

func (s *search) addToQueue(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "search.addToQueue"

	s.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID
	media := s.mediaSelectedStorage.Get(chatId)

	segm, err := s.sch.AddToQueue(ctx, chatId, media)
	if err != nil {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
			return
		}
		return
	}

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    chatId,
		MessageID: s.msgIdStorage.Get(chatId),
		Text:      s.successMsg(segm.Start, segm.Start.Add(segm.StopCut-segm.BeginCut)),
	}); err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
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

	id := s.mediaPageStorage.Get(chatId)
	res := s.mediaResultsStorage.Get(chatId)

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

func (s *search) deleteMedia(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "search.deleteMedia"

	s.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   s.msgIdStorage.Get(chatId),
		Text:        ctr.LibSearchDeleteSubmit,
		ReplyMarkup: s.submitDeleteMarkup(),
	}); err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (s *search) deleteSubmit(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "search.deleteSubmit"

	s.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	if err := s.lib.DeleteMedia(ctx, chatId, s.mediaSelectedStorage.Get(chatId)); err != nil {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	s.mediaResultsStorage.Set(chatId, []localModels.MediaConfig{})

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    chatId,
		MessageID: s.msgIdStorage.Get(chatId),
		Text:      ctr.LibSearchDeleteSuccess,
	}); err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (s *search) deleteReject(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "search.closedUpdateMedia"

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	id := s.mediaPageStorage.Get(chatId)
	res := s.mediaResultsStorage.Get(chatId)

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

func (s *search) successMsg(start, stop time.Time) string {
	return fmt.Sprintf("Добавлено в расписание с %s по %s.", start.Format("06-01-02 15:04:05"), stop.Format("15:04:05"))
}

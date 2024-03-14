package search

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/split"
)

func (s *Search) update(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	var msg string

	switch s.router.GetState(update.CallbackQuery.Data) {
	case "name-author":
		s.targetUpdateStorage.Set(chatId, "name-author")
		msg = ctr.LibSearchAskNameAuthor
	case "genre":
		s.targetUpdateStorage.Set(chatId, "genre")
		msg = ctr.LibSearchAskGenre
	case "format":
		opt := s.searchStorage.Get(chatId)
		switch opt.format {
		case formatSong:
			opt.format = formatPodcast
			opt.playlists = nil
		case formatPodcast:
			opt.format = formatSong
			opt.podcasts = nil
		}
		s.searchStorage.Set(chatId, opt)
		id := s.msgIdStorage.Get(chatId)

		if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatId,
			MessageID:   id,
			Text:        s.filterRepr(opt),
			ReplyMarkup: s.mainMenuMarkup(opt),
			ParseMode:   models.ParseModeHTML,
		}); err != nil {
			s.onError(err)
		}
		return
	case "podcast-playlist":
		opt := s.searchStorage.Get(chatId)
		switch opt.format {
		case formatSong:
			msg = ctr.LibSearchAskPlaylist

		case formatPodcast:
			msg = ctr.LibSearchAskPodcast
		}
		s.targetUpdateStorage.Set(chatId, "podcast-playlist")
	case "lang":
		s.targetUpdateStorage.Set(chatId, "lang")
		msg = ctr.LibSearchAskLang
	case "mood":
		s.targetUpdateStorage.Set(chatId, "mood")
		msg = ctr.LibSearchAskMood
	case "reset":
		opt := s.searchStorage.Get(chatId)
		opt.genres = nil
		opt.playlists = nil
		opt.podcasts = nil
		opt.languages = nil
		opt.moods = nil
		s.searchStorage.Set(chatId, opt)
		msg = s.filterRepr(opt)

		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    chatId,
			Text:      msg,
			ParseMode: models.ParseModeHTML,
		}); err != nil {
			s.onError(err)
		}
		return
	}

	s.session.Redirect(chatId, s.router.Path(cmdGetData))

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   s.msgIdStorage.Get(chatId),
		Text:        msg,
		ReplyMarkup: s.getSettingDataMarkup(),
	}); err != nil {
		s.onError(err)
	}
}

func (s *Search) getData(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID

	opt := s.searchStorage.Get(chatId)
	msg := update.Message.Text

	if msg == "" {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.LibUploadErrEmptyMsg,
		}); err != nil {
			s.onError(err)
		}
		return
	}

	switch s.targetUpdateStorage.Get(chatId) {
	case "name-author":
		opt.nameAuthor = msg
	case "genre":
		opt.genres = split.SplitMsg(msg)
	case "podcast-playlist":
		switch opt.format {
		case formatSong:
			opt.playlists = split.SplitMsg(msg)
		case formatPodcast:
			opt.podcasts = split.SplitMsg(msg)
		}
	case "lang":
		opt.languages = split.SplitMsg(msg)
	case "mood":
		opt.moods = split.SplitMsg(msg)
	default:
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			s.onError(err)
		}
		return
	}

	s.session.Redirect(chatId, ctr.NullStatus)
	s.targetUpdateStorage.Del(chatId)
	s.searchStorage.Set(chatId, opt)

	if _, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    chatId,
		MessageID: update.Message.ID,
	}); err != nil {
		s.onError(err)
	}

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   s.msgIdStorage.Get(chatId),
		Text:        s.filterRepr(opt),
		ReplyMarkup: s.mainMenuMarkup(opt),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		s.onError(err)
	}
}

func (s *Search) cancelSlider(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	opt := s.searchStorage.Get(chatId)

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   s.msgIdStorage.Get(chatId),
		Text:        s.filterRepr(opt),
		ReplyMarkup: s.mainMenuMarkup(opt),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		s.onError(err)
	}
}

package search

import (
	"context"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (s *Search) updateState(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID

	msg := update.Message.Text
	target := s.targetUpdateStorage.Get(userId)

	opt := s.searchStorage.Get(userId)

	switch target {
	case cmdNameAuthor:
		if msg == "" {
			if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatId,
				Text:   ctr.LibSearchErrNameAuthorEmpty,
			}); err != nil {
				s.onError(err)
			}
			return
		}
		opt.nameAuthor = msg
	case cmdFormat:
		switch opt.format {
		case formatAll:
			opt.format = formatSong
		case formatSong:
			opt.format = formatPodcast
		case formatPodcast:
			opt.format = formatAll
		}
	case cmdPlaylist:
		opt.playlists = splitMsg(msg)
	case cmdGenre:
		opt.genres = splitMsg(msg)
	case cmdLanguage:
		opt.languages = splitMsg(msg)
	case cmdMood:
		opt.moods = splitMsg(msg)
	default:
		s.mediaResults.Del(userId)
	}

	s.targetUpdateStorage.Set(userId, "")
	s.searchStorage.Set(userId, opt)
	s.session.Redirect(userId, ctr.NullStatus)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        s.filterRepr(userId),
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: s.formatSelectMarkup(),
	}); err != nil {
		s.onError(err)
	}
}

func (s *Search) nameAuthor(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	s.targetUpdateStorage.Set(userId, cmdNameAuthor)
	s.session.Redirect(userId, s.router.Path(cmdUpdateOption))

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   ctr.LibSearchAskNameAuthror,
	}); err != nil {
		s.onError(err)
	}
}

func (s *Search) format(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	opt := s.searchStorage.Get(userId)
	switch opt.format {
	case formatAll:
		opt.format = formatSong
	case formatSong:
		opt.format = formatPodcast
	case formatPodcast:
		opt.format = formatAll
	}
	s.searchStorage.Set(userId, opt)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        ctr.LibSearchAskFormat,
		ReplyMarkup: s.formatSelectMarkup(),
	}); err != nil {
		s.onError(err)
	}
}

func (s *Search) formatAll(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	opt := s.searchStorage.Get(userId)
	opt.format = formatAll
	s.searchStorage.Set(userId, opt)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        s.filterRepr(userId),
		ReplyMarkup: s.mainMenuMarkup(),
	}); err != nil {
		s.onError(err)
	}
}

func (s *Search) formatPodcast(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	opt := s.searchStorage.Get(userId)
	opt.format = formatPodcast
	s.searchStorage.Set(userId, opt)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        s.filterRepr(userId),
		ReplyMarkup: s.mainMenuMarkup(),
	}); err != nil {
		s.onError(err)
	}
}

func (s *Search) formatSong(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	opt := s.searchStorage.Get(userId)
	opt.format = formatSong
	s.searchStorage.Set(userId, opt)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        s.filterRepr(userId),
		ReplyMarkup: s.mainMenuMarkup(),
	}); err != nil {
		s.onError(err)
	}
}

func (s *Search) playlist(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	s.targetUpdateStorage.Set(userId, cmdPlaylist)
	s.session.Redirect(userId, s.router.Path(cmdUpdateOption))

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   ctr.LibSearchAskPlaylist,
	}); err != nil {
		s.onError(err)
	}
}

func (s *Search) genre(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	s.targetUpdateStorage.Set(userId, cmdGenre)
	s.session.Redirect(userId, s.router.Path(cmdUpdateOption))

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   ctr.LibSearchAskGenre,
	}); err != nil {
		s.onError(err)
	}
}

func (s *Search) language(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	s.targetUpdateStorage.Set(userId, cmdGenre)
	s.session.Redirect(userId, s.router.Path(cmdUpdateOption))

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   ctr.LibSearchAskLanguage,
	}); err != nil {
		s.onError(err)
	}
}

func (s *Search) mood(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	s.targetUpdateStorage.Set(userId, cmdMood)
	s.session.Redirect(userId, s.router.Path(cmdUpdateOption))

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   ctr.LibSearchAskMood,
	}); err != nil {
		s.onError(err)
	}
}

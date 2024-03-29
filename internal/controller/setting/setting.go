package setting

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/split"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/storage"
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
)

const (
	cmdBase          ctr.Command = ""
	cmdUpdateSetting ctr.Command = "update-setting"
	cmdGetData       ctr.Command = "get-data"
	cmdOpenCheckBox  ctr.Command = "check-open"
	cmdCheckBtn      ctr.Command = "check-btn"
	cmdCancelSetting ctr.Command = "cancel"
	cmdSubmit        ctr.Command = "submit"
	cmdClose         ctr.Command = "close"

	cmdNoOp ctr.Command = "no-op"
)

type setting struct {
	ctr.CallbackAnswerer

	router   *ctr.Router
	session  ctr.Session
	onSelect ctr.OnSelectHandler
	onCancel ctr.OnSelectHandler
	onError  bot.ErrorsHandler

	initialConfigStorage storage.Storage[localModels.MediaConfig]
	mediaConfigStorage   storage.Storage[localModels.MediaConfig]
	msgIdStorage         storage.Storage[int]
	targetStorage        storage.Storage[string]
}

type OnSelect func()

func Register(
	router *ctr.Router,
	session ctr.Session,
	onSelect ctr.OnSelectHandler,
	onCancel ctr.OnSelectHandler,
	onError bot.ErrorsHandler,
	mediaConfigStorage storage.Storage[localModels.MediaConfig],
	msgIdStorage storage.Storage[int],
) {
	s := &setting{
		router:   router,
		session:  session,
		onSelect: onSelect,
		onCancel: onCancel,
		onError:  onError,

		mediaConfigStorage:   mediaConfigStorage,
		msgIdStorage:         msgIdStorage,
		initialConfigStorage: storage.New[localModels.MediaConfig](),
		targetStorage:        storage.New[string](),
	}

	router.RegisterCallback(cmdBase, s.init)

	router.RegisterCallbackPrefix(cmdUpdateSetting, s.updateSettings)
	router.RegisterHandler(cmdGetData, s.getSettingNewData)
	router.RegisterCallbackPrefix(cmdOpenCheckBox, s.openCheckBox)
	router.RegisterCallbackPrefix(cmdCheckBtn, s.getCheckedBtn)
	router.RegisterCallback(cmdCancelSetting, s.cancelSubTask)

	router.RegisterCallback(cmdSubmit, s.submit)
	router.RegisterCallback(cmdClose, s.close)

	router.RegisterCallback(cmdNoOp, s.nullHandler)
}

func (s *setting) init(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "setting.init"

	s.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	conf := s.mediaConfigStorage.Get(chatId)
	s.initialConfigStorage.Set(chatId, conf)

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   s.msgIdStorage.Get(chatId),
		Text:        conf.String(),
		ReplyMarkup: s.MainSettingsMarkup(conf),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (s *setting) updateSettings(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.updateSettings"

	s.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	var msg string

	switch s.router.GetState(update.CallbackQuery.Data) {
	case "name":
		s.targetStorage.Set(chatId, "name")
		msg = ctr.LibSearchUpdateAskName
	case "author":
		s.targetStorage.Set(chatId, "author")
		msg = ctr.LibSearchUpdateAskAuthor
	case "format":
		conf := s.mediaConfigStorage.Get(chatId)
		switch conf.Format {
		case localModels.Song:
			conf.Format = localModels.Podcast
			conf.Playlists = nil
		case localModels.Podcast:
			conf.Format = localModels.Jingle
			conf.Podcasts = nil
		case localModels.Jingle:
			conf.Format = localModels.Song
		}
		s.mediaConfigStorage.Set(chatId, conf)

		if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatId,
			MessageID:   s.msgIdStorage.Get(chatId),
			Text:        conf.String(),
			ReplyMarkup: s.MainSettingsMarkup(conf),
			ParseMode:   models.ParseModeHTML,
		}); err != nil {
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	case "podcast-playlist":
		conf := s.mediaConfigStorage.Get(chatId)
		var state string
		switch conf.Format {
		case localModels.Song:
			state = "playlists"
			msg = ctr.LibSearchUpdateAskPlaylist
		case localModels.Podcast:
			state = "podcasts"
			msg = ctr.LibSearchUpdateAskPodcast
		}
		s.targetStorage.Set(chatId, state)
	case "reset":
		conf := s.initialConfigStorage.Get(chatId)
		s.mediaConfigStorage.Set(chatId, conf)
		msg = conf.String()

		if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatId,
			MessageID:   s.msgIdStorage.Get(chatId),
			Text:        msg,
			ReplyMarkup: s.MainSettingsMarkup(conf),
			ParseMode:   models.ParseModeHTML,
		}); err != nil {
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	default:
		return
	}

	s.session.Redirect(chatId, s.router.Path(cmdGetData))

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   s.msgIdStorage.Get(chatId),
		Text:        msg,
		ReplyMarkup: s.getSettingDataMarkup(),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (s *setting) getSettingNewData(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.getSettingNewData"

	chatId := update.Message.Chat.ID

	conf := s.mediaConfigStorage.Get(chatId)
	msg := update.Message.Text

	if msg == "" {
		if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    chatId,
			MessageID: s.msgIdStorage.Get(chatId),
			Text:      ctr.LibSearchUpdateErrEmptyMsg,
		}); err != nil {
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	switch s.targetStorage.Get(chatId) {
	case "name":
		conf.Name = msg
	case "author":
		conf.Author = msg
	case "podcasts":
		conf.Podcasts = split.SplitMsg(msg)
	case "playlists":
		conf.Playlists = split.SplitMsg(msg)
	}

	s.session.Redirect(chatId, ctr.NullStatus)
	s.targetStorage.Del(chatId)
	s.mediaConfigStorage.Set(chatId, conf)

	if _, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    chatId,
		MessageID: update.Message.ID,
	}); err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   s.msgIdStorage.Get(chatId),
		Text:        conf.String(),
		ReplyMarkup: s.MainSettingsMarkup(conf),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (s *setting) openCheckBox(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.openCheckBox"

	s.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID
	callback := s.router.GetState(update.CallbackQuery.Data)

	var (
		msg    string
		markup models.InlineKeyboardMarkup
	)

	conf := s.mediaConfigStorage.Get(chatId)

	switch callback {
	case "genre":
		msg = ctr.LibUploadAskGenre
		markup = s.genreChooseMarkup(conf)
	case "mood":
		msg = ctr.LibUploadAskMood
		markup = s.moodChooseMarkup(conf)
	case "lang":
		msg = ctr.LibUploadAskLang
		markup = s.langChooseMarkup(conf)
	}

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   s.msgIdStorage.Get(chatId),
		Text:        msg,
		ReplyMarkup: markup,
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (s *setting) getCheckedBtn(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.getCheckedBtn"

	s.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID
	callback := s.router.GetState(update.CallbackQuery.Data)

	// tagType in ('genre', 'mood', 'lang')
	// data is tag id (from 1 to GenreNumber, MoodNumber, LangNumber)
	tagType, data, found := strings.Cut(callback, "-")
	if !found {
		s.onError(fmt.Errorf("%s [%d]: invalid callback data \"%s\"", op, chatId, callback))
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}
	id, err := strconv.Atoi(data)
	if err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	conf := s.mediaConfigStorage.Get(chatId)

	var (
		msg    string
		markup models.InlineKeyboardMarkup
	)

	switch tagType {
	case "genre":
		conf.Genres[id-1] = !conf.Genres[id-1]
		msg = ctr.LibSearchUpdateAskGenre
		markup = s.genreChooseMarkup(conf)
	case "mood":
		conf.Moods[id-1] = !conf.Moods[id-1]
		msg = ctr.LibSearchUpdateAskMood
		markup = s.moodChooseMarkup(conf)
	case "lang":
		conf.Languages[id-1] = !conf.Languages[id-1]
		msg = ctr.LibSearchUpdateAskLang
		markup = s.langChooseMarkup(conf)
	}

	s.mediaConfigStorage.Set(chatId, conf)

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   s.msgIdStorage.Get(chatId),
		Text:        msg,
		ReplyMarkup: markup,
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (s *setting) cancelSubTask(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.cancelSubTask"

	s.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	conf := s.mediaConfigStorage.Get(chatId)

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   s.msgIdStorage.Get(chatId),
		Text:        conf.String(),
		ReplyMarkup: s.MainSettingsMarkup(conf),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (s *setting) submit(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.CallbackAnswer(ctx, b, update.CallbackQuery)

	s.onSelect(ctx, b, update.CallbackQuery.Message)
}

func (s *setting) close(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.CallbackAnswer(ctx, b, update.CallbackQuery)

	s.onCancel(ctx, b, update.CallbackQuery.Message)
}

func (s *setting) nullHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.CallbackAnswer(ctx, b, update.CallbackQuery)
}

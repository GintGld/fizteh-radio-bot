package upload

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/split"
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (u *upload) updateSettings(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.updateSettings"

	u.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	var msg string

	switch u.router.GetState(update.CallbackQuery.Data) {
	case "name":
		u.settingTargetStorage.Set(chatId, "name")
		msg = ctr.LibUploadAskName
	case "author":
		u.settingTargetStorage.Set(chatId, "author")
		msg = ctr.LibUploadAskAuthor
	case "genre":
		u.settingTargetStorage.Set(chatId, "genre")
		msg = ctr.LibUploadAskGenre
	case "format":
		conf := u.mediaConfigStorage.Get(chatId)
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
		u.mediaConfigStorage.Set(chatId, conf)

		if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatId,
			MessageID:   u.msgIdStorage.Get(chatId),
			Text:        conf.String(),
			ReplyMarkup: u.mediaConfMarkup(conf),
			ParseMode:   models.ParseModeHTML,
		}); err != nil {
			u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	case "podcast-playlist":
		conf := u.mediaConfigStorage.Get(chatId)
		var state string
		switch conf.Format {
		case localModels.Song:
			state = "playlists"
			msg = ctr.LibUploadAskPlaylist
		case localModels.Podcast:
			state = "podcasts"
			msg = ctr.LibUploadAskPodcast
		}
		u.settingTargetStorage.Set(chatId, state)
	case "lang":
		u.settingTargetStorage.Set(chatId, "lang")
		msg = ctr.LibUploadAskLang
	case "mood":
		u.settingTargetStorage.Set(chatId, "mood")
		msg = ctr.LibUploadAskMood
	case "reset":
		conf := u.mediaConfigStorage.Get(chatId)
		conf.Genres = [localModels.GenreNumber]bool{}
		conf.Playlists = nil
		conf.Podcasts = nil
		conf.Languages = [localModels.LangNumber]bool{}
		conf.Moods = [localModels.MoodNumber]bool{}
		u.mediaConfigStorage.Set(chatId, conf)
		msg = conf.String()

		if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatId,
			MessageID:   u.msgIdStorage.Get(chatId),
			Text:        msg,
			ReplyMarkup: u.mediaConfMarkup(conf),
			ParseMode:   models.ParseModeHTML,
		}); err != nil {
			u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	default:
		return
	}

	u.session.Redirect(chatId, u.router.Path(cmdGetData))

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   u.msgIdStorage.Get(chatId),
		Text:        msg,
		ReplyMarkup: u.getSettingDataMarkup(),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (u *upload) getSettingNewData(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.getSettingNewData"

	chatId := update.Message.Chat.ID

	conf := u.mediaConfigStorage.Get(chatId)
	msg := update.Message.Text

	if msg == "" {
		if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    chatId,
			MessageID: u.msgIdStorage.Get(chatId),
			Text:      ctr.LibUploadErrEmptyMsg,
		}); err != nil {
			u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	switch u.settingTargetStorage.Get(chatId) {
	case "name":
		conf.Name = msg
	case "author":
		conf.Author = msg
	case "podcasts":
		conf.Podcasts = split.SplitMsg(msg)
	case "playlists":
		conf.Playlists = split.SplitMsg(msg)
	}

	u.session.Redirect(chatId, ctr.NullStatus)
	u.settingTargetStorage.Del(chatId)
	u.mediaConfigStorage.Set(chatId, conf)

	if _, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    chatId,
		MessageID: update.Message.ID,
	}); err != nil {
		u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   u.msgIdStorage.Get(chatId),
		Text:        conf.String(),
		ReplyMarkup: u.mediaConfMarkup(conf),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (u *upload) openCheckBox(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.openCheckBox"

	u.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID
	callback := u.router.GetState(update.CallbackQuery.Data)

	var (
		msg    string
		markup models.InlineKeyboardMarkup
	)

	conf := u.mediaConfigStorage.Get(chatId)

	switch callback {
	case "genre":
		msg = ctr.LibSearchUpdateAskGenre
		markup = u.genreChooseMarkup(conf)
	case "mood":
		msg = ctr.LibSearchUpdateAskMood
		markup = u.moodChooseMarkup(conf)
	case "lang":
		msg = ctr.LibSearchUpdateAskLang
		markup = u.langChooseMarkup(conf)
	}

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   u.msgIdStorage.Get(chatId),
		Text:        msg,
		ReplyMarkup: markup,
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (u *upload) getCheckedBtn(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.getCheckedBtn"

	u.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID
	callback := u.router.GetState(update.CallbackQuery.Data)

	// tagType in ('genre', 'mood', 'lang')
	// data is tag id (from 1 to GenreNumber, MoodNumber, LangNumber)
	tagType, data, found := strings.Cut(callback, "-")
	if !found {
		u.onError(fmt.Errorf("%s [%d]: invalid callback data \"%s\"", op, chatId, callback))
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}
	id, err := strconv.Atoi(data)
	if err != nil {
		u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	conf := u.mediaConfigStorage.Get(chatId)

	var (
		msg    string
		markup models.InlineKeyboardMarkup
	)

	switch tagType {
	case "genre":
		conf.Genres[id-1] = !conf.Genres[id-1]
		msg = ctr.LibSearchUpdateAskGenre
		markup = u.genreChooseMarkup(conf)
	case "mood":
		conf.Moods[id-1] = !conf.Moods[id-1]
		msg = ctr.LibSearchUpdateAskMood
		markup = u.moodChooseMarkup(conf)
	case "lang":
		conf.Languages[id-1] = !conf.Languages[id-1]
		msg = ctr.LibSearchUpdateAskLang
		markup = u.langChooseMarkup(conf)
	}

	u.mediaConfigStorage.Set(chatId, conf)

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   u.msgIdStorage.Get(chatId),
		Text:        msg,
		ReplyMarkup: markup,
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (u *upload) cancelSubTask(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.cancelSubTask"

	u.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	conf := u.mediaConfigStorage.Get(chatId)

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   u.msgIdStorage.Get(chatId),
		Text:        conf.String(),
		ReplyMarkup: u.mediaConfMarkup(conf),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

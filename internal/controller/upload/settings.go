package upload

import (
	"context"
	"fmt"
	"strings"
	"time"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/split"
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (u *upload) updateSettings(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.updateSettings"

	u.callbackAnswer(ctx, b, update.CallbackQuery)

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
			conf.Format = localModels.Song
			conf.Podcasts = nil
		}
		u.mediaConfigStorage.Set(chatId, conf)

		if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatId,
			MessageID:   u.msgIdStorage.Get(chatId),
			Text:        u.mediaConfRepr(conf),
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
		conf.Genres = nil
		conf.Playlists = nil
		conf.Podcasts = nil
		conf.Languages = nil
		conf.Moods = nil
		u.mediaConfigStorage.Set(chatId, conf)
		msg = u.mediaConfRepr(conf)

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
	case "genre":
		conf.Genres = split.SplitMsg(msg)
	case "podcasts":
		conf.Podcasts = split.SplitMsg(msg)
	case "playlists":
		conf.Playlists = split.SplitMsg(msg)
	case "lang":
		conf.Languages = split.SplitMsg(msg)
	case "mood":
		conf.Moods = split.SplitMsg(msg)
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
		Text:        u.mediaConfRepr(conf),
		ReplyMarkup: u.mediaConfMarkup(conf),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (u *upload) cancelSubTask(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.cancelSubTask"

	u.callbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	conf := u.mediaConfigStorage.Get(chatId)

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   u.msgIdStorage.Get(chatId),
		Text:        u.mediaConfRepr(conf),
		ReplyMarkup: u.mediaConfMarkup(conf),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (u *upload) mediaConfRepr(conf localModels.MediaConfig) string {
	var b strings.Builder

	b.WriteString("<b>Композиция</b>\n")
	b.WriteString(fmt.Sprintf("<b>Название:</b> %s\n", conf.Name))
	b.WriteString(fmt.Sprintf("<b>Автор:</b> %s\n", conf.Author))
	b.WriteString(fmt.Sprintf("<b>Длительность:</b> %s\n", conf.Duration.Round(time.Second).String()))

	if len(conf.Podcasts) > 0 {
		b.WriteString(fmt.Sprintf("<b>Подкасты:</b> %s\n", strings.Join(conf.Podcasts, ", ")))
	}
	if len(conf.Playlists) > 0 {
		b.WriteString(fmt.Sprintf("<b>Плейлисты:</b> %s\n", strings.Join(conf.Playlists, ", ")))
	}
	if len(conf.Genres) > 0 {
		b.WriteString(fmt.Sprintf("<b>Жанры:</b> %s\n", strings.Join(conf.Genres, ", ")))
	}
	if len(conf.Languages) > 0 {
		b.WriteString(fmt.Sprintf("<b>Языки:</b> %s\n", strings.Join(conf.Languages, ", ")))
	}
	if len(conf.Moods) > 0 {
		b.WriteString(fmt.Sprintf("<b>Настроение:</b> %s\n", strings.Join(conf.Moods, ", ")))
	}

	return b.String()
}

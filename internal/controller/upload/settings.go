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
	u.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	var msg string

	switch u.router.GetState(update.CallbackQuery.Data) {
	case "name":
		u.settingTargetStorage.Set(userId, "name")
		msg = ctr.LibUploadAskName
	case "author":
		u.settingTargetStorage.Set(userId, "author")
		msg = ctr.LibUploadAskAuthor
	case "genre":
		u.settingTargetStorage.Set(userId, "genre")
		msg = ctr.LibUploadAskGenre
	case "format":
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      chatId,
			Text:        ctr.LibUploadAskFormat,
			ReplyMarkup: u.formatMarkup(),
		}); err != nil {
			u.onError(err)
		}
		return
	case "podcast-playlist":
		conf := u.mediaConfigStorage.Get(userId)
		var state string
		switch conf.Format {
		case localModels.Song:
			state = "song"
			msg = ctr.LibUploadAskPlaylist
		case localModels.Podcast:
			state = "podcast"
			msg = ctr.LibUploadAskPodcast
		}
		u.settingTargetStorage.Set(userId, state)
	case "lang":
		u.settingTargetStorage.Set(userId, "lang")
		msg = ctr.LibUploadAskLang
	case "mood":
		u.settingTargetStorage.Set(userId, "mood")
		msg = ctr.LibUploadAskMood
	case "reset":
		conf := u.mediaConfigStorage.Get(userId)
		conf.Genres = nil
		conf.Playlists = nil
		conf.Podcasts = nil
		conf.Languages = nil
		conf.Moods = nil
		u.mediaConfigStorage.Set(userId, conf)
		msg = u.mediaConfRepr(conf)

		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   msg,
		}); err != nil {
			u.onError(err)
		}
		return
	}

	u.session.Redirect(userId, u.router.Path(cmdGetData))

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        msg,
		ReplyMarkup: u.getSettingDataMarkup(),
	}); err != nil {
		u.onError(err)
	}
}

func (u *upload) getSettingNewData(ctx context.Context, b *bot.Bot, update *models.Update) {
	u.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID

	conf := u.mediaConfigStorage.Get(userId)
	msg := update.Message.Text

	if msg == "" {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.LibUploadErrEmptyMsg,
		}); err != nil {
			u.onError(err)
		}
		return
	}

	switch u.settingTargetStorage.Get(userId) {
	case "name":
		conf.Name = msg
	case "author":
		conf.Author = msg
	case "genre":
		conf.Genres = split.SplitMsg(msg)
	case "podcast-playlist":
		switch conf.Format {
		case localModels.Song:
			conf.Playlists = split.SplitMsg(msg)
		case localModels.Podcast:
			conf.Podcasts = split.SplitMsg(msg)
		}
	case "lang":
		conf.Languages = split.SplitMsg(msg)
	case "mood":
		conf.Moods = split.SplitMsg(msg)
	}

	u.session.Redirect(userId, ctr.NullStatus)
	u.settingTargetStorage.Del(userId)
	u.mediaConfigStorage.Set(userId, conf)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        u.mediaConfRepr(conf),
		ReplyMarkup: u.mediaConfMarkup(conf),
	}); err != nil {
		u.onError(err)
	}
}

func (u *upload) updateFormat(ctx context.Context, b *bot.Bot, update *models.Update) {
	u.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	conf := u.mediaConfigStorage.Get(userId)

	switch u.router.GetState(update.CallbackQuery.Data) {
	case "song":
		conf.Format = localModels.Song
	case "podcast":
		conf.Format = localModels.Podcast
	}

	u.mediaConfigStorage.Set(userId, conf)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        u.mediaConfRepr(conf),
		ReplyMarkup: u.mediaConfMarkup(conf),
	}); err != nil {
		u.onError(err)
	}
}

func (u *upload) cancelSubTask(ctx context.Context, b *bot.Bot, update *models.Update) {
	u.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	conf := u.mediaConfigStorage.Get(userId)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        u.mediaConfRepr(conf),
		ReplyMarkup: u.mediaConfMarkup(conf),
	}); err != nil {
		u.onError(err)
	}
}

func (u *upload) mediaConfRepr(conf localModels.MediaConfig) string {
	var b strings.Builder

	b.WriteString("*Композиция*")
	b.WriteString(fmt.Sprintf("*Название:* %s\n", conf.Name))
	b.WriteString(fmt.Sprintf("*Автор:* %s\n", conf.Author))
	b.WriteString(fmt.Sprintf("*Длительность:* %s\n", conf.Duration.Round(time.Second).String()))

	if len(conf.Podcasts) > 0 {
		b.WriteString(fmt.Sprintf("*Подкасты*: %s\n", strings.Join(conf.Podcasts, ", ")))
	}
	if len(conf.Playlists) > 0 {
		b.WriteString(fmt.Sprintf("*Плейлисты*: %s\n", strings.Join(conf.Playlists, ", ")))
	}
	if len(conf.Genres) > 0 {
		b.WriteString(fmt.Sprintf("*Жанры*: %s\n", strings.Join(conf.Genres, ", ")))
	}
	if len(conf.Languages) > 0 {
		b.WriteString(fmt.Sprintf("*Языки*: %s\n", strings.Join(conf.Languages, ", ")))
	}
	if len(conf.Moods) > 0 {
		b.WriteString(fmt.Sprintf("*Настроение*: %s\n", strings.Join(conf.Moods, ", ")))
	}

	return b.String()
}

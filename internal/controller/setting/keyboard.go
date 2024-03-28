package setting

import (
	"github.com/go-telegram/bot/models"

	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
)

const (
	butMsgName     = "Название"
	butMsgAuthor   = "Автор"
	butMsgGenre    = "Жанр"
	butMsgAlbum    = "Альбом"
	butMsgPlaylist = "Плейлисты"
	butMsgPodcast  = "Подкаст"
	butMsgPodcasts = "Подкасты"
	butMsgLang     = "Язык"
	butMsgMood     = "Настроение"
	butMsgSong     = "Песня"
	butMsgJingle   = "Джингл"
	butMsgReset    = "Сбросить"

	butMsgSubmit = "Сохранить"
	butMsgCancel = "Назад"
)

func (s *setting) MainSettingsMarkup(conf localModels.MediaConfig) models.InlineKeyboardMarkup {
	var target, form, formCallback string
	switch conf.Format {
	case localModels.Podcast:
		target = butMsgPodcast
		form = butMsgPodcasts
		formCallback = "podcast-playlist"
	case localModels.Song:
		target = butMsgSong
		form = butMsgPlaylist
		formCallback = "podcast-playlist"
	case localModels.Jingle:
		target = butMsgJingle
		form = "\t"
		formCallback = ""
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgName, CallbackData: s.router.PathPrefixState(cmdUpdateSetting, "name")},
				{Text: butMsgAuthor, CallbackData: s.router.PathPrefixState(cmdUpdateSetting, "author")},
			},
			{
				{Text: target, CallbackData: s.router.PathPrefixState(cmdUpdateSetting, "format")},
				{Text: form, CallbackData: s.router.PathPrefixState(cmdUpdateSetting, formCallback)},
			},
			{
				{Text: butMsgAlbum, CallbackData: s.router.PathPrefixState(cmdUpdateSetting, "album")},
				{Text: butMsgGenre, CallbackData: s.router.PathPrefixState(cmdUpdateSetting, "genre")},
			},
			{
				{Text: butMsgLang, CallbackData: s.router.PathPrefixState(cmdUpdateSetting, "lang")},
				{Text: butMsgMood, CallbackData: s.router.PathPrefixState(cmdUpdateSetting, "mood")},
			},
			{
				{Text: butMsgReset, CallbackData: s.router.PathPrefixState(cmdUpdateSetting, "reset")},
				{Text: butMsgSubmit, CallbackData: s.router.Path(cmdSubmit)},
			},
			{
				{Text: butMsgCancel, CallbackData: s.router.Path(cmdClose)},
			},
		},
	}
}

func (s *setting) getSettingDataMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgCancel, CallbackData: s.router.Path(cmdCancelSetting)},
			},
		},
	}
}

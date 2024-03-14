package upload

import (
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
	"github.com/go-telegram/bot/models"
)

const (
	butMsgManual = "Файл"
	butMsgLink   = "Ссылка"

	butMsgName     = "Название"
	butMsgAuthor   = "Автор"
	butMsgGenre    = "Жанр"
	butMsgFormat   = "Формат"
	butMsgPlaylist = "Плейлисты"
	butMsgPodcast  = "Подкаст"
	butMsgPodcasts = "Подкасты"
	butMsgLang     = "Язык"
	butMsgMood     = "Настроение"
	butMsgSong     = "Песня"
	butMsgReset    = "Сбросить"

	butMsgSubmit = "Загрузить"
	butMsgCancel = "Назад"
)

func (u *upload) mainMenuMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgManual, CallbackData: u.router.Path(cmdManual)},
				{Text: butMsgLink, CallbackData: u.router.Path(cmdLink)},
			},
		},
	}
}

func (u *upload) mediaConfMarkup(conf localModels.MediaConfig) models.InlineKeyboardMarkup {
	var target, form string
	switch conf.Format {
	case localModels.Podcast:
		target = butMsgPodcast
		form = butMsgPodcasts
	case localModels.Song:
		target = butMsgSong
		form = butMsgPlaylist
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgName, CallbackData: u.router.PathPrefixState(cmdSettings, "name")},
				{Text: butMsgAuthor, CallbackData: u.router.PathPrefixState(cmdSettings, "author")},
			},
			{
				{Text: target, CallbackData: u.router.PathPrefixState(cmdSettings, "format")},
				{Text: form, CallbackData: u.router.PathPrefixState(cmdSettings, "podcast-playlist")},
			},
			{
				{Text: butMsgGenre, CallbackData: u.router.PathPrefixState(cmdSettings, "genre")},
				{Text: butMsgLang, CallbackData: u.router.PathPrefixState(cmdSettings, "lang")},
			},
			{
				{Text: butMsgMood, CallbackData: u.router.PathPrefixState(cmdSettings, "mood")},
				{Text: butMsgReset, CallbackData: u.router.PathPrefixState(cmdSettings, "reset")},
			},
			{
				{Text: butMsgSubmit, CallbackData: u.router.PathPrefixState(cmdSubmit, "manual")},
				{Text: butMsgCancel, CallbackData: u.router.Path(cmdCancel)},
			},
		},
	}
}

func (u *upload) getSettingDataMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgCancel, CallbackData: u.router.Path(cmdCancelSetting)},
			},
		},
	}
}

func (u *upload) playlistMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgSubmit, CallbackData: u.router.PathPrefixState(cmdSubmit, "link")},
				{Text: butMsgCancel, CallbackData: u.router.Path(cmdCancel)},
			},
		},
	}
}

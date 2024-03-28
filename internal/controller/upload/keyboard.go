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
	butMsgPlaylist = "Плейлисты"
	butMsgPodcast  = "Подкаст"
	butMsgPodcasts = "Подкасты"
	butMsgLang     = "Язык"
	butMsgMood     = "Настроение"
	butMsgSong     = "Песня"
	butMsgJingle   = "Джингл"
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
				{Text: butMsgName, CallbackData: u.router.PathPrefixState(cmdSettings, "name")},
				{Text: butMsgAuthor, CallbackData: u.router.PathPrefixState(cmdSettings, "author")},
			},
			{
				{Text: target, CallbackData: u.router.PathPrefixState(cmdSettings, "format")},
				{Text: form, CallbackData: u.router.PathPrefixState(cmdSettings, formCallback)},
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
				{Text: butMsgSubmit, CallbackData: u.router.Path(cmdSubmit)},
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

func (u *upload) cancelMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgCancel, CallbackData: u.router.Path(cmdCancel)},
			},
		},
	}
}

func (u *upload) askUploadMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgSubmit, CallbackData: u.router.Path(cmdSubmit)},
				{Text: butMsgCancel, CallbackData: u.router.Path(cmdCancel)},
			},
		},
	}
}

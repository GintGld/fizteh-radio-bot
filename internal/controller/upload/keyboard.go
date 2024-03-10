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
	butMsgPlaylist = "Плейлист"
	butMsgPodcast  = "Подкаст"
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
			{
				{Text: butMsgCancel, CallbackData: u.router.Path(cmdBack)},
			},
		},
	}
}

func (u *upload) mediaConfMarkup(conf localModels.MediaConfig) models.InlineKeyboardMarkup {
	var form string
	switch conf.Format {
	case localModels.Podcast:
		form = butMsgPodcast
	case localModels.Song:
		form = butMsgPlaylist
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgName, CallbackData: u.router.PathPrefixState(cmdSettings, "name")},
				{Text: butMsgAuthor, CallbackData: u.router.PathPrefixState(cmdSettings, "author")},
			},
			{
				{Text: butMsgFormat, CallbackData: u.router.PathPrefixState(cmdSettings, "format")},
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

func (u *upload) formatMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgSong, CallbackData: u.router.PathPrefixState(cmdFormat, "song")},
				{Text: butMsgPodcast, CallbackData: u.router.PathPrefixState(cmdFormat, "podcast")},
			},
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

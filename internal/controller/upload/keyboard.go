package upload

import (
	"fmt"

	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/slice"
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
	"github.com/go-telegram/bot/models"
)

const (
	butMsgManual = "Файл"
	butMsgLink   = "Ссылка"

	butMsgName       = "Название"
	butMsgAuthor     = "Автор"
	butMsgGenre      = "Жанр"
	butMsgPlaylist   = "Плейлисты"
	butMsgPodcast    = "Подкаст"
	butMsgPodcasts   = "Подкасты"
	butMsgLang       = "Язык"
	butMsgMood       = "Настроение"
	butMsgSong       = "Песня"
	butMsgJingle     = "Джингл"
	butMsgReset      = "Сбросить"
	butMsgChecked    = "☑️"
	butMsgNotChecked = "✖️"

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
				{Text: butMsgGenre, CallbackData: u.router.PathPrefixState(cmdOpenCheckBox, "genre")},
				{Text: butMsgLang, CallbackData: u.router.PathPrefixState(cmdOpenCheckBox, "lang")},
			},
			{
				{Text: butMsgMood, CallbackData: u.router.PathPrefixState(cmdOpenCheckBox, "mood")},
				{Text: butMsgReset, CallbackData: u.router.PathPrefixState(cmdSettings, "reset")},
			},
			{
				{Text: butMsgSubmit, CallbackData: u.router.Path(cmdSubmit)},
				{Text: butMsgCancel, CallbackData: u.router.Path(cmdCancel)},
			},
		},
	}
}

func (u *upload) genreChooseMarkup(conf localModels.MediaConfig) models.InlineKeyboardMarkup {
	var msg string

	const rowLen = 2

	rows := make([][]models.InlineKeyboardButton, 0, localModels.GenreNumber/rowLen)
	row := make([]models.InlineKeyboardButton, 0, rowLen)

	for _, g := range localModels.GenresAvail {
		msg = g.String()
		if conf.Genres[g.Id-1] {
			msg += butMsgChecked
		} else {
			msg += butMsgNotChecked
		}
		callback := fmt.Sprintf("genre-%d", g.Id)

		if len(row) == rowLen {
			rows = append(rows, row)
			row = make([]models.InlineKeyboardButton, 0, rowLen)
		}
		row = append(row, models.InlineKeyboardButton{
			Text:         msg,
			CallbackData: u.router.PathPrefixState(cmdCheckBtn, callback),
		})

	}

	if len(row) > 0 && len(row) < rowLen {
		row = append(row, slice.Repeat(
			models.InlineKeyboardButton{
				Text:         "\t",
				CallbackData: u.router.Path(cmdNoOp),
			},
			rowLen-len(row),
		)...)
	}

	rows = append(rows, row, []models.InlineKeyboardButton{{
		Text:         butMsgCancel,
		CallbackData: u.router.Path(cmdCancelSetting),
	}})

	return models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func (u *upload) moodChooseMarkup(conf localModels.MediaConfig) models.InlineKeyboardMarkup {
	var msg string

	const rowLen = 2

	rows := make([][]models.InlineKeyboardButton, 0, localModels.MoodNumber/rowLen)
	row := make([]models.InlineKeyboardButton, 0, rowLen)

	for _, m := range localModels.MoodsAvail {
		msg = m.String()
		if conf.Moods[m.Id-1] {
			msg += butMsgChecked
		} else {
			msg += butMsgNotChecked
		}
		callback := fmt.Sprintf("mood-%d", m.Id)

		if len(row) == rowLen {
			rows = append(rows, row)
			row = make([]models.InlineKeyboardButton, 0, rowLen)
		}
		row = append(row, models.InlineKeyboardButton{
			Text:         msg,
			CallbackData: u.router.PathPrefixState(cmdCheckBtn, callback),
		})
	}

	if len(row) > 0 && len(row) < rowLen {
		row = append(row, slice.Repeat(
			models.InlineKeyboardButton{
				Text:         "\t",
				CallbackData: u.router.Path(cmdNoOp),
			},
			rowLen-len(row),
		)...)
	}

	rows = append(rows, row, []models.InlineKeyboardButton{{
		Text:         butMsgCancel,
		CallbackData: u.router.Path(cmdCancelSetting),
	}})

	return models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func (u *upload) langChooseMarkup(conf localModels.MediaConfig) models.InlineKeyboardMarkup {
	var msg string

	const rowLen = 3

	rows := make([][]models.InlineKeyboardButton, 0, localModels.MoodNumber/rowLen)
	row := make([]models.InlineKeyboardButton, 0, rowLen)

	for _, l := range localModels.LangsAvail {
		msg = l.String()
		if conf.Languages[l.Id-1] {
			msg += butMsgChecked
		} else {
			msg += butMsgNotChecked
		}
		callback := fmt.Sprintf("lang-%d", l.Id)

		if len(row) == rowLen {
			rows = append(rows, row)
			row = make([]models.InlineKeyboardButton, 0, rowLen)
		}
		row = append(row, models.InlineKeyboardButton{
			Text:         msg,
			CallbackData: u.router.PathPrefixState(cmdCheckBtn, callback),
		})
	}

	if len(row) > 0 && len(row) < rowLen {
		row = append(row, slice.Repeat(
			models.InlineKeyboardButton{
				Text:         "\t",
				CallbackData: u.router.Path(cmdNoOp),
			},
			rowLen-len(row),
		)...)
	}

	rows = append(rows, row, []models.InlineKeyboardButton{{
		Text:         butMsgCancel,
		CallbackData: u.router.Path(cmdCancelSetting),
	}})

	return models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
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

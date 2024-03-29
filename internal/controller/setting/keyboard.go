package setting

import (
	"fmt"

	"github.com/go-telegram/bot/models"

	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/slice"
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
)

const (
	butMsgName       = "Название"
	butMsgAuthor     = "Автор"
	butMsgGenre      = "Жанр"
	butMsgAlbum      = "Альбом"
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
				{Text: butMsgGenre, CallbackData: s.router.PathPrefixState(cmdOpenCheckBox, "genre")},
			},
			{
				{Text: butMsgLang, CallbackData: s.router.PathPrefixState(cmdOpenCheckBox, "lang")},
				{Text: butMsgMood, CallbackData: s.router.PathPrefixState(cmdOpenCheckBox, "mood")},
			},
			{
				{Text: butMsgReset, CallbackData: s.router.PathPrefixState(cmdUpdateSetting, "reset")},
				{Text: butMsgCancel, CallbackData: s.router.Path(cmdClose)},
			},
			{
				{Text: butMsgSubmit, CallbackData: s.router.Path(cmdSubmit)},
			},
		},
	}
}

func (s *setting) genreChooseMarkup(conf localModels.MediaConfig) models.InlineKeyboardMarkup {
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
			CallbackData: s.router.PathPrefixState(cmdCheckBtn, callback),
		})
	}

	if len(row) < rowLen {
		row = append(row, slice.Repeat(
			models.InlineKeyboardButton{
				Text:         "\t",
				CallbackData: s.router.Path(cmdNoOp),
			},
			rowLen-len(row),
		)...)
	}

	rows = append(rows, row, []models.InlineKeyboardButton{{
		Text:         butMsgCancel,
		CallbackData: s.router.Path(cmdCancelSetting),
	}})

	return models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func (s *setting) moodChooseMarkup(conf localModels.MediaConfig) models.InlineKeyboardMarkup {
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
			CallbackData: s.router.PathPrefixState(cmdCheckBtn, callback),
		})
	}

	if len(row) < rowLen {
		row = append(row, slice.Repeat(
			models.InlineKeyboardButton{
				Text:         "\t",
				CallbackData: s.router.Path(cmdNoOp),
			},
			rowLen-len(row),
		)...)
	}

	rows = append(rows, row, []models.InlineKeyboardButton{{
		Text:         butMsgCancel,
		CallbackData: s.router.Path(cmdCancelSetting),
	}})

	return models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func (s *setting) langChooseMarkup(conf localModels.MediaConfig) models.InlineKeyboardMarkup {
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
			CallbackData: s.router.PathPrefixState(cmdCheckBtn, callback),
		})
	}

	if len(row) < localModels.LangNumber {
		row = append(row, slice.Repeat(
			models.InlineKeyboardButton{
				Text:         "\t",
				CallbackData: s.router.Path(cmdNoOp),
			},
			rowLen-len(row),
		)...)
	}

	rows = append(rows, row, []models.InlineKeyboardButton{{
		Text:         butMsgCancel,
		CallbackData: s.router.Path(cmdCancelSetting),
	}})

	return models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
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

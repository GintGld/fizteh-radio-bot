package autodj

import (
	"fmt"

	"github.com/go-telegram/bot/models"

	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/slice"
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
)

const (
	butMsgAlbum      = "Альбомы"
	butMsgGenre      = "Жанры"
	butMsgPlaylist   = "Плейлисты"
	butMsgLanguage   = "Языки"
	butMsgMood       = "Настроения"
	butMsgReset      = "Сбросить"
	butMsgUpdate     = "Обновить"
	butMsgCancel     = "Назад"
	butMsgSend       = "Обновить настройки"
	butMsgStart      = "Запустить"
	butMsgStop       = "Остановить"
	butMsgChecked    = "☑️"
	butMsgNotChecked = "✖️"
)

func (a *autodj) mainMenuMarkup(conf localModels.AutoDJInfo) models.InlineKeyboardMarkup {
	playing := butMsgStart
	if conf.IsPlaying {
		playing = butMsgStop
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgGenre, CallbackData: a.router.PathPrefixState(cmdOpenCheckBox, "genre")},
				{Text: butMsgPlaylist, CallbackData: a.router.PathPrefixState(cmdUpdate, "playlist")},
			},
			{
				{Text: butMsgLanguage, CallbackData: a.router.PathPrefixState(cmdOpenCheckBox, "lang")},
				{Text: butMsgMood, CallbackData: a.router.PathPrefixState(cmdOpenCheckBox, "mood")},
			},
			{
				{Text: butMsgAlbum, CallbackData: a.router.PathPrefixState(cmdUpdate, "album")},
				{Text: butMsgReset, CallbackData: a.router.Path(cmdReset)},
			},
			{
				{Text: butMsgSend, CallbackData: a.router.Path(cmdSend)},
				{Text: playing, CallbackData: a.router.Path(cmdStartStop)},
			},
		},
	}
}

func (a *autodj) genreChooseMarkup(conf localModels.AutoDJInfo) models.InlineKeyboardMarkup {
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
			CallbackData: a.router.PathPrefixState(cmdCheckBtn, callback),
		})

	}

	if len(row) > 0 && len(row) < rowLen {
		row = append(row, slice.Repeat(
			models.InlineKeyboardButton{
				Text:         "\t",
				CallbackData: a.router.Path(cmdNoOp),
			},
			rowLen-len(row),
		)...)
	}

	rows = append(rows, row, []models.InlineKeyboardButton{{
		Text:         butMsgCancel,
		CallbackData: a.router.Path(cmdCloseSubtask),
	}})

	return models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func (a *autodj) moodChooseMarkup(conf localModels.AutoDJInfo) models.InlineKeyboardMarkup {
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
			CallbackData: a.router.PathPrefixState(cmdCheckBtn, callback),
		})
	}

	if len(row) > 0 && len(row) < rowLen {
		row = append(row, slice.Repeat(
			models.InlineKeyboardButton{
				Text:         "\t",
				CallbackData: a.router.Path(cmdNoOp),
			},
			rowLen-len(row),
		)...)
	}

	rows = append(rows, row, []models.InlineKeyboardButton{{
		Text:         butMsgCancel,
		CallbackData: a.router.Path(cmdCloseSubtask),
	}})

	return models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func (a *autodj) langChooseMarkup(conf localModels.AutoDJInfo) models.InlineKeyboardMarkup {
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
			CallbackData: a.router.PathPrefixState(cmdCheckBtn, callback),
		})
	}

	if len(row) > 0 && len(row) < rowLen {
		row = append(row, slice.Repeat(
			models.InlineKeyboardButton{
				Text:         "\t",
				CallbackData: a.router.Path(cmdNoOp),
			},
			rowLen-len(row),
		)...)
	}

	rows = append(rows, row, []models.InlineKeyboardButton{{
		Text:         butMsgCancel,
		CallbackData: a.router.Path(cmdCloseSubtask),
	}})

	return models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

package library

import (
	"github.com/go-telegram/bot/models"
)

const (
	buttonSearch = "Искать"
	buttonUpload = "Загрузить"
)

func (l *library) mainMenuMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{{
			models.InlineKeyboardButton{
				Text:         buttonSearch,
				CallbackData: l.router.Path(cmdSearch),
			},
			models.InlineKeyboardButton{
				Text:         buttonUpload,
				CallbackData: l.router.Path(cmdUpload),
			},
		}},
	}
}

package schedule

import "github.com/go-telegram/bot/models"

const (
	buttonLook = "Текущее расписание"
	buttonDj   = "Dj"
)

func (s *schedule) mainMenuMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{{
			models.InlineKeyboardButton{
				Text:         buttonLook,
				CallbackData: cmdLook,
			},
			models.InlineKeyboardButton{
				Text:         buttonDj,
				CallbackData: cmdDj,
			},
		}},
	}
}

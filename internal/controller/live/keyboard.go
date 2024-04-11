package live

import "github.com/go-telegram/bot/models"

const (
	BtnStart = "Начать эфир"
	BtnStop  = "Остановить"
	BtnYes   = "Да"
	BtnNo    = "Нет"
)

func (l *live) StartMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: BtnStart, CallbackData: l.router.Path(cmdStart)}},
		},
	}
}

func (l *live) StopMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: BtnStop, CallbackData: l.router.Path(cmdStop)}},
		},
	}
}

func (l *live) SubmitStopMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: BtnYes, CallbackData: l.router.Path(cmdStopSubmit)},
				{Text: BtnNo, CallbackData: l.router.Path(cmdStopReject)},
			},
		},
	}
}

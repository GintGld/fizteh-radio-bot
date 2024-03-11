package look

import (
	"strconv"

	"github.com/go-telegram/bot/models"
)

const (
	butMsgPrev = "\u00AB"
	butMsgNext = "\u00BB"

	butMsgUpdate = "Обновить"
	butMsgCancel = "Назад"
)

func (l *look) mainMenuMarkup(page, maxPage int) models.InlineKeyboardMarkup {
	var (
		butLeft = models.InlineKeyboardButton{
			Text:         butMsgPrev,
			CallbackData: l.router.PathPrefixState(cmdNewPage, strconv.Itoa(page-1)),
		}
		butRight = models.InlineKeyboardButton{
			Text:         butMsgNext,
			CallbackData: l.router.PathPrefixState(cmdNewPage, strconv.Itoa(page+1)),
		}
	)
	if page == 1 {
		butLeft = models.InlineKeyboardButton{Text: "", CallbackData: l.router.Path(cmdNoOp)}
	}
	if page == maxPage {
		butRight = models.InlineKeyboardButton{Text: "", CallbackData: l.router.Path(cmdNoOp)}
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				butLeft,
				butRight,
			},
			{
				{Text: butMsgUpdate, CallbackData: l.router.Path(cmdUpdate)},
				{Text: butMsgCancel, CallbackData: l.router.Path(cmdBack)},
			},
		},
	}
}

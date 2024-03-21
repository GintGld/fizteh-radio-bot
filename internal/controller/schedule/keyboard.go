package schedule

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

func (s *schedule) mainMenuMarkup(page, maxPage int) models.InlineKeyboardMarkup {
	var (
		butLeft = models.InlineKeyboardButton{
			Text:         butMsgPrev,
			CallbackData: s.router.PathPrefixState(cmdNewPage, strconv.Itoa(page-1)),
		}
		butRight = models.InlineKeyboardButton{
			Text:         butMsgNext,
			CallbackData: s.router.PathPrefixState(cmdNewPage, strconv.Itoa(page+1)),
		}
	)
	if page == 1 {
		butLeft = models.InlineKeyboardButton{Text: "\t", CallbackData: s.router.Path(cmdNoOp)}
	}
	if page == maxPage {
		butRight = models.InlineKeyboardButton{Text: "\t", CallbackData: s.router.Path(cmdNoOp)}
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				butLeft,
				butRight,
			},
			{
				{Text: butMsgUpdate, CallbackData: s.router.Path(cmdUpdate)},
			},
		},
	}
}

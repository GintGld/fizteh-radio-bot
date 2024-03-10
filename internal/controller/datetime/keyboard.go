package datetime

import "github.com/go-telegram/bot/models"

const (
	butMsgCancel = "Назад"
)

func (p *picker) timePickMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "00:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "00:00")},
				{Text: "00:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "00:30")},
				{Text: "01:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "01:00")},
				{Text: "01:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "01:30")},
				{Text: "02:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "02:00")},
				{Text: "02:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "02:30")},
			},
			{
				{Text: "03:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "03:00")},
				{Text: "03:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "03:30")},
				{Text: "04:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "04:00")},
				{Text: "04:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "04:30")},
				{Text: "05:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "05:00")},
				{Text: "05:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "05:30")},
			},
			{
				{Text: "06:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "06:00")},
				{Text: "06:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "06:30")},
				{Text: "07:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "07:00")},
				{Text: "07:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "07:30")},
				{Text: "08:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "08:00")},
				{Text: "08:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "08:30")},
			},
			{
				{Text: "09:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "09:00")},
				{Text: "09:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "09:30")},
				{Text: "10:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "10:00")},
				{Text: "10:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "10:30")},
				{Text: "11:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "11:00")},
				{Text: "11:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "11:30")},
			},
			{
				{Text: "12:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "12:00")},
				{Text: "12:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "12:30")},
				{Text: "13:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "13:00")},
				{Text: "13:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "13:30")},
				{Text: "14:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "14:00")},
				{Text: "14:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "14:30")},
			},
			{
				{Text: "15:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "15:00")},
				{Text: "15:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "15:30")},
				{Text: "16:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "16:00")},
				{Text: "16:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "16:30")},
				{Text: "17:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "17:00")},
				{Text: "17:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "17:30")},
			},
			{
				{Text: "18:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "18:00")},
				{Text: "18:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "18:30")},
				{Text: "19:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "19:00")},
				{Text: "19:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "19:30")},
				{Text: "20:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "20:00")},
				{Text: "20:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "20:30")},
			},
			{
				{Text: "21:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "21:00")},
				{Text: "21:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "21:30")},
				{Text: "22:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "22:00")},
				{Text: "22:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "22:30")},
				{Text: "23:00", CallbackData: p.router.PathPrefixState(cmdSubmit, "23:00")},
				{Text: "23:30", CallbackData: p.router.PathPrefixState(cmdSubmit, "23:30")},
			},
			{
				{Text: butMsgCancel, CallbackData: p.router.Path(cmdCancelTime)},
			},
		},
	}
}

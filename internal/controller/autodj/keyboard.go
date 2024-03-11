package autodj

import "github.com/go-telegram/bot/models"

const (
	butMsgGenre    = "Жанры"
	butMsgPlaylist = "Плейлисты"
	butMsgLanguage = "Языки"
	butMsgMood     = "Настроения"
	butMsgReset    = "Сбросить"
	butMsgSubmit   = "Обновить"
	butMsgCancel   = "Назад"
)

func (a *autodj) mainMenuMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgGenre, CallbackData: a.router.PathPrefixState(cmdUpdate, "genre")},
				{Text: butMsgPlaylist, CallbackData: a.router.PathPrefixState(cmdUpdate, "playlist")},
			},
			{
				{Text: butMsgLanguage, CallbackData: a.router.PathPrefixState(cmdUpdate, "lang")},
				{Text: butMsgMood, CallbackData: a.router.PathPrefixState(cmdUpdate, "mood")},
			},
			{
				{Text: butMsgReset, CallbackData: a.router.Path(cmdReset)},
				{Text: "", CallbackData: a.router.Path(cmdNoOp)},
			},
			{
				{Text: butMsgSubmit, CallbackData: a.router.Path(cmdSubmit)},
				{Text: butMsgCancel, CallbackData: a.router.Path(cmdBack)},
			},
		},
	}
}

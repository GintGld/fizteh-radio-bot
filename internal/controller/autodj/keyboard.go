package autodj

import (
	"github.com/go-telegram/bot/models"

	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
)

const (
	butMsgGenre    = "Жанры"
	butMsgPlaylist = "Плейлисты"
	butMsgLanguage = "Языки"
	butMsgMood     = "Настроения"
	butMsgReset    = "Сбросить"
	butMsgUpdate   = "Обновить"
	butMsgCancel   = "Назад"
	butMsgSend     = "Обновить настройки"
	butMsgStart    = "Запустить"
	butMsgStop     = "Остановить"
)

func (a *autodj) mainMenuMarkup(conf localModels.AutoDJInfo) models.InlineKeyboardMarkup {
	playing := butMsgStart
	if conf.IsPlaying {
		playing = butMsgStop
	}

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
				{Text: butMsgUpdate, CallbackData: a.router.Path(cmdGetCurrConf)},
				{Text: butMsgReset, CallbackData: a.router.Path(cmdReset)},
			},
			{
				{Text: butMsgSend, CallbackData: a.router.Path(cmdSend)},
				{Text: playing, CallbackData: a.router.Path(cmdStartStop)},
			},
		},
	}
}

package search

import (
	"fmt"

	"github.com/go-telegram/bot/models"
)

const (
	butMsgNameOrAuthor  = "Название/автор"
	butMsgFormat        = "Песни/подкасты"
	butMsgFormatAll     = "Все"
	butMsgFormatPodcast = "Подкасты"
	butMsgFormatSong    = "Песни"
	butMsgPlaylist      = "Плейлисты"
	butMsgGenre         = "Жанры"
	butMsgLanguage      = "Язык"
	butMsgMood          = "Настроение"

	butMsgSubmit = "Искать"
	butMsgCancel = "Назад"
)

func (s *Search) mainMenuMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgNameOrAuthor, CallbackData: s.router.Path(cmdNameAuthor)},
				{Text: butMsgFormat, CallbackData: s.router.Path(cmdFormat)},
			},
			{
				{Text: butMsgPlaylist, CallbackData: s.router.Path(cmdPlaylist)},
				{Text: butMsgGenre, CallbackData: s.router.Path(cmdGenre)},
			},
			{
				{Text: butMsgLanguage, CallbackData: s.router.Path(cmdLanguage)},
				{Text: butMsgMood, CallbackData: s.router.Path(cmdMood)},
			},
			{
				{Text: butMsgSubmit, CallbackData: s.router.Path(cmdSubmit)},
				{Text: butMsgCancel, CallbackData: s.router.Path(cmdBase)},
			},
		},
	}
}

func (s *Search) formatSelectMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgFormatAll, CallbackData: s.router.Path(cmdFormatAll)},
				{Text: butMsgFormatPodcast, CallbackData: s.router.Path(cmdFormatPodcast)},
				{Text: butMsgFormatSong, CallbackData: s.router.Path(cmdFormatSong)},
			},
		},
	}
}

func (s *Search) mediaSliderMarkup(id int, maxId int) models.InlineKeyboardMarkup {
	var (
		butLeft = models.InlineKeyboardButton{
			Text:         "\u00AB",
			CallbackData: s.router.Path(cmdPrevSlide),
		}
		butRight = models.InlineKeyboardButton{
			Text:         "\u00BB",
			CallbackData: s.router.Path(cmdNextSlide),
		}
	)
	if id == 1 {
		butLeft.Text = ""
		butLeft.CallbackData = s.router.Path(cmdNoOp)
	}
	if id == maxId {
		butRight.Text = ""
		butRight.CallbackData = s.router.Path(cmdNoOp)
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				butLeft,
				{Text: fmt.Sprintf("%d/%d", id, maxId), CallbackData: s.router.Path(cmdNoOp)},
				butRight,
			},
			{
				{Text: "В расписание", CallbackData: s.router.Path(cmdNoOp)},  // TODO
				{Text: "Редактировать", CallbackData: s.router.Path(cmdNoOp)}, //TODO
			},
			{
				{Text: "Назад", CallbackData: s.router.Path(cmdCloseSlider)},
			},
		},
	}
}

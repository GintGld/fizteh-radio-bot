package search

import (
	"fmt"

	"github.com/go-telegram/bot/models"
)

const (
	butMsgNameOrAuthor = "Название/автор"
	butMsgPodcast      = "Подкасты"
	butMsgSong         = "Песни"
	butMsgPlaylist     = "Плейлисты"
	butMsgPodcasts     = "Подкасты"
	butMsgGenre        = "Жанры"
	butMsgLanguage     = "Язык"
	butMsgMood         = "Настроение"

	butMsgAddToSch = "Запланировать"
	butMsgEdit     = "Редактировать"
	butMsgPlayNext = "Играть следующим"

	butMsgSubmit = "Искать"
	butMsgCancel = "Назад"
)

func (s *Search) mainMenuMarkup(opt searchOption) models.InlineKeyboardMarkup {
	msgFormat := butMsgSong
	msgFormatSelect := butMsgPlaylist
	if opt.format == formatPodcast {
		msgFormat = butMsgPodcast
		msgFormatSelect = butMsgPodcasts
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgNameOrAuthor, CallbackData: s.router.PathPrefixState(cmdUpdate, "name-author")},
				{Text: msgFormat, CallbackData: s.router.PathPrefixState(cmdUpdate, "format")},
			},
			{
				{Text: msgFormatSelect, CallbackData: s.router.PathPrefixState(cmdUpdate, "podcast-playlist")},
				{Text: butMsgGenre, CallbackData: s.router.PathPrefixState(cmdUpdate, "genre")},
			},
			{
				{Text: butMsgLanguage, CallbackData: s.router.PathPrefixState(cmdUpdate, "lang")},
				{Text: butMsgMood, CallbackData: s.router.PathPrefixState(cmdUpdate, "mood")},
			},
			{
				{Text: butMsgSubmit, CallbackData: s.router.Path(cmdSubmit)},
			},
		},
	}
}

func (s *Search) getSettingDataMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgCancel, CallbackData: s.router.Path(cmdCloseSlider)},
			},
		},
	}
}

func (s *Search) mediaSliderMarkup(id int, maxId int) models.InlineKeyboardMarkup {
	var (
		butLeft = models.InlineKeyboardButton{
			Text:         "\u00AB",
			CallbackData: s.router.PathPrefixState(cmdUpdateSlide, "prev"),
		}
		butRight = models.InlineKeyboardButton{
			Text:         "\u00BB",
			CallbackData: s.router.PathPrefixState(cmdUpdateSlide, "next"),
		}
	)
	if id == 1 {
		butLeft.Text = "\t"
		butLeft.CallbackData = s.router.Path(cmdNoOp)
	}
	if id == maxId {
		butRight.Text = "\t"
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
				{Text: butMsgAddToSch, CallbackData: s.router.Path(cmdSelectMedia)},
				{Text: butMsgPlayNext, CallbackData: s.router.Path(cmdNoOp)}, // TODO
			},
			{
				{Text: butMsgEdit, CallbackData: s.router.Path(cmdNoOp)}, // TODO
				{Text: butMsgCancel, CallbackData: s.router.Path(cmdCloseSlider)},
			},
		},
	}
}

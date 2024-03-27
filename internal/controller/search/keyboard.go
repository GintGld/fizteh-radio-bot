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
	butMsgJingle       = "Джинглы"
	butMsgPodcasts     = "Подкасты"
	butMsgGenre        = "Жанры"
	butMsgLanguage     = "Язык"
	butMsgMood         = "Настроение"
	butMsgReset        = "Сбросить"

	butMsgAddToSch = "Запланировать"
	butMsgEdit     = "Редактировать"
	butMsgPlayNext = "Играть следующим"
	butMsgDelete   = "Удалить"

	butMsgSubmit = "Искать"
	butMsgCancel = "Назад"
)

func (s *search) mainMenuMarkup(opt searchOption) models.InlineKeyboardMarkup {
	var msgFormat, msgFormatSelect, msgCallback string
	switch opt.format {
	case formatSong:
		msgFormat = butMsgSong
		msgFormatSelect = butMsgPlaylist
		msgCallback = "podcast-playlist"
	case formatPodcast:
		msgFormat = butMsgPodcast
		msgFormatSelect = butMsgPodcasts
		msgCallback = "podcast-playlist"
	case formatJingle:
		msgFormat = butMsgJingle
		msgFormatSelect = "\t"
		msgCallback = ""
	}

	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgNameOrAuthor, CallbackData: s.router.PathPrefixState(cmdUpdate, "name-author")},
				{Text: msgFormat, CallbackData: s.router.PathPrefixState(cmdUpdate, "format")},
			},
			{
				{Text: msgFormatSelect, CallbackData: s.router.PathPrefixState(cmdUpdate, msgCallback)},
				{Text: butMsgGenre, CallbackData: s.router.PathPrefixState(cmdUpdate, "genre")},
			},
			{
				{Text: butMsgLanguage, CallbackData: s.router.PathPrefixState(cmdUpdate, "lang")},
				{Text: butMsgMood, CallbackData: s.router.PathPrefixState(cmdUpdate, "mood")},
			},
			{
				{Text: butMsgSubmit, CallbackData: s.router.Path(cmdSubmit)},
				{Text: butMsgReset, CallbackData: s.router.Path(cmdReset)},
			},
		},
	}
}

func (s *search) getSettingDataMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: butMsgCancel, CallbackData: s.router.Path(cmdCloseSlider)},
			},
		},
	}
}

func (s *search) mediaSliderMarkup(id int, maxId int) models.InlineKeyboardMarkup {
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
				{Text: butMsgEdit, CallbackData: s.router.Path(cmdUpdateMediaInfo)},
				{Text: butMsgDelete, CallbackData: s.router.Path(cmdDeleteMedia)},
			},
			{
				{Text: butMsgCancel, CallbackData: s.router.Path(cmdCloseSlider)},
			},
		},
	}
}

func (s *search) submitDeleteMarkup() models.InlineKeyboardMarkup {
	return models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Да", CallbackData: s.router.Path(cmdDeleteSubmit)},
				{Text: "Нет", CallbackData: s.router.Path(cmdDeleteReject)},
			},
		},
	}
}

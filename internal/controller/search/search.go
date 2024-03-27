package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/GintGld/fizteh-radio-bot/internal/controller/datetime"
	"github.com/GintGld/fizteh-radio-bot/internal/controller/setting"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/storage"
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
)

const (
	// basic
	cmdBase   ctr.Command = ""
	cmdSubmit ctr.Command = "text"

	// options
	cmdUpdate  ctr.Command = "update"
	cmdGetData ctr.Command = "get-data"
	cmdFormat  ctr.Command = "format"
	cmdReset   ctr.Command = "reset"

	// media slider
	cmdUpdateSlide     ctr.Command = "update-slide"
	cmdCloseSlider     ctr.Command = "cancel"
	cmdSelectMedia     ctr.Command = "select"
	cmdUpdateMediaInfo ctr.Command = "update-media"

	// filler
	cmdNoOp ctr.Command = "no-op"
)

type search struct {
	ctr.CallbackAnswerer

	router  *ctr.Router
	auth    Auth
	lib     Library
	session ctr.Session
	onError bot.ErrorsHandler

	searchStorage        storage.Storage[searchOption]
	targetUpdateStorage  storage.Storage[string]
	mediaPage            storage.Storage[int]
	mediaResults         storage.Storage[[]localModels.MediaConfig]
	mediaSelectedStorage storage.Storage[localModels.MediaConfig]
	msgIdStorage         storage.Storage[int]
}

type Auth interface {
	IsKnown(ctx context.Context, id int64) bool
}

type Library interface {
	Search(ctx context.Context, id int64, filter localModels.MediaFilter) ([]localModels.MediaConfig, error)
	UpdateMedia(ctx context.Context, id int64, mediaConf localModels.MediaConfig) error
}

type searchOption struct {
	nameAuthor string
	format     searchFormat
	playlists  []string
	podcasts   []string
	genres     []string
	languages  []string
	moods      []string
}

type searchFormat int

const (
	formatSong searchFormat = iota
	formatPodcast
	formatJingle
)

func (opt searchOption) ToFilter() localModels.MediaFilter {
	tags := make([]string, 0)
	tags = append(tags, opt.format.String())
	tags = append(tags, opt.playlists...)
	tags = append(tags, opt.genres...)
	tags = append(tags, opt.languages...)
	tags = append(tags, opt.moods...)

	return localModels.MediaFilter{
		Name:       opt.nameAuthor,
		Author:     opt.nameAuthor,
		Tags:       tags,
		MaxRespLen: 20,
	}
}

func (sOpt searchFormat) String() string {
	switch sOpt {
	case formatSong:
		return "song"
	case formatPodcast:
		return "podcast"
	case formatJingle:
		return "jingle"
	default:
		return ""
	}
}

func (sOpt searchFormat) Repr() string {
	switch sOpt {
	case formatSong:
		return "песня"
	case formatPodcast:
		return "подкаст"
	case formatJingle:
		return "джингл"
	default:
		return ""
	}
}

func Register(
	router *ctr.Router,
	auth Auth,
	lib Library,
	scheduleAdd datetime.ScheduleAdd,
	session ctr.Session,
	onError bot.ErrorsHandler,
) {
	s := &search{
		router:  router,
		auth:    auth,
		lib:     lib,
		session: session,
		onError: onError,

		searchStorage:        storage.New[searchOption](),
		targetUpdateStorage:  storage.New[string](),
		mediaPage:            storage.New[int](),
		mediaResults:         storage.New[[]localModels.MediaConfig](),
		mediaSelectedStorage: storage.New[localModels.MediaConfig](),
		msgIdStorage:         storage.New[int](),
	}

	// main menu
	// router.RegisterCallback(cmdBase, s.init)
	router.RegisterCommand(s.init)

	// call searcher and show result
	router.RegisterCallback(cmdSubmit, s.submit)

	// option updates
	router.RegisterCallbackPrefix(cmdUpdate, s.update)
	router.RegisterHandler(cmdGetData, s.getData)
	router.RegisterCallback(cmdReset, s.reset)

	// media slider
	router.RegisterCallbackPrefix(cmdUpdateSlide, s.updateSlide)
	router.RegisterCallback(cmdCloseSlider, s.cancelSlider)

	// selector for schedule modify
	datetime.Register(
		router.With(cmdSelectMedia),
		scheduleAdd,
		session,
		s.canceledDateTimeSelector,
		onError,
		s.msgIdStorage,
		s.mediaSelectedStorage,
	)

	// selector for updating media info
	setting.Register(
		router.With(cmdUpdateMediaInfo),
		session,
		s.updateMedia,
		s.closedUpdateMedia,
		onError,
		s.mediaSelectedStorage,
		s.msgIdStorage,
	)

	// null handler to answer callbacks for empty buttons
	router.RegisterCallback(cmdNoOp, s.nullHandler)
}

func (s *search) init(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "search.init"

	chatId := update.Message.Chat.ID

	if !s.auth.IsKnown(ctx, chatId) {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrUnknown,
		}); err != nil {
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	opt := s.searchStorage.Get(chatId)

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        s.filterRepr(opt),
		ReplyMarkup: s.mainMenuMarkup(opt),
		ParseMode:   models.ParseModeHTML,
	})
	if err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}

	s.msgIdStorage.Set(chatId, msg.ID)
}

func (s *search) reset(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "search.submit"

	s.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	s.searchStorage.Set(chatId, searchOption{})

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   s.msgIdStorage.Get(chatId),
		Text:        s.filterRepr(searchOption{}),
		ReplyMarkup: s.mainMenuMarkup(searchOption{}),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

// submit gets text to search,
// search option and returns
// slider for found media.
func (s *search) submit(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "search.submit"

	s.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	opt := s.searchStorage.Get(chatId)

	res, err := s.lib.Search(ctx, chatId, opt.ToFilter())
	// TODO enhance errors
	if err != nil {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	if len(res) == 0 {
		if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatId,
			MessageID:   s.msgIdStorage.Get(chatId),
			Text:        ctr.LibSearchErrEmptyRes,
			ReplyMarkup: s.getSettingDataMarkup(),
		}); err != nil {
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	s.mediaPage.Set(chatId, 1)
	s.mediaResults.Set(chatId, res)
	s.mediaSelectedStorage.Set(chatId, res[0])

	msg, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   s.msgIdStorage.Get(chatId),
		Text:        res[0].String(),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: s.mediaSliderMarkup(1, len(res)),
	})

	if err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}

	s.msgIdStorage.Set(chatId, msg.ID)
}

func (s *search) nullHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.CallbackAnswer(ctx, b, update.CallbackQuery)
}

// filterRepr returns formated filter info.
func (s *search) filterRepr(opt searchOption) string {
	var b strings.Builder

	b.WriteString("<b>Настройки поиска:</b>\n")
	if opt.nameAuthor != "" {
		b.WriteString(fmt.Sprintf("<b>Название/автор:</b> %s\n", opt.nameAuthor))
	}
	b.WriteString(fmt.Sprintf("<b>Формат:</b> %s\n", opt.format.Repr()))
	if len(opt.playlists) > 0 {
		b.WriteString(fmt.Sprintf("<b>Плейлисты:</b> %s\n", strings.Join(opt.playlists, ", ")))
	}
	if len(opt.podcasts) > 0 {
		b.WriteString(fmt.Sprintf("<b>Подкасты:</b> %s\n", strings.Join(opt.podcasts, ", ")))
	}
	if len(opt.genres) > 0 {
		b.WriteString(fmt.Sprintf("<b>Жанры:</b> %s\n", strings.Join(opt.genres, ", ")))
	}
	if len(opt.languages) > 0 {
		b.WriteString(fmt.Sprintf("<b>Языки:</b> %s\n", strings.Join(opt.languages, ", ")))
	}
	if len(opt.moods) > 0 {
		b.WriteString(fmt.Sprintf("<b>Настроения:</b> %s", strings.Join(opt.moods, ", ")))
	}

	return b.String()
}

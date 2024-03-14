package search

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/GintGld/fizteh-radio-bot/internal/controller/datetime"
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

	// media slider
	cmdUpdateSlide ctr.Command = "update-slide"
	cmdCloseSlider ctr.Command = "cancel"
	cmdSelectMedia ctr.Command = "select"

	// filler
	cmdNoOp ctr.Command = "no-op"
)

type Search struct {
	router  *ctr.Router
	auth    Auth
	search  LibrarySearch
	session ctr.Session
	onError bot.ErrorsHandler

	searchStorage       storage.Storage[searchOption]
	targetUpdateStorage storage.Storage[string]
	mediaPage           storage.Storage[int]
	mediaResults        storage.Storage[[]localModels.Media]
	mediaSelected       storage.Storage[localModels.Media]
	msgIdStorage        storage.Storage[int]
}

type Auth interface {
	IsKnown(id int64) bool
}

type LibrarySearch interface {
	Search(id int64, filter localModels.MediaFilter) ([]localModels.Media, error)
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
)

func (sOpt searchFormat) String() string {
	switch sOpt {
	case formatSong:
		return "песня"
	case formatPodcast:
		return "подкаст"
	default:
		return ""
	}
}

func Register(
	router *ctr.Router,
	auth Auth,
	search LibrarySearch,
	scheduleAdd datetime.ScheduleAdd,
	session ctr.Session,
	onError bot.ErrorsHandler,
) {
	s := &Search{
		router:  router,
		auth:    auth,
		search:  search,
		session: session,
		onError: onError,

		searchStorage:       storage.New[searchOption](),
		targetUpdateStorage: storage.New[string](),
		mediaPage:           storage.New[int](),
		mediaResults:        storage.New[[]localModels.Media](),
		mediaSelected:       storage.New[localModels.Media](),
		msgIdStorage:        storage.New[int](),
	}

	// main menu
	// router.RegisterCallback(cmdBase, s.init)
	router.RegisterCommand(s.init)

	// call searcher and show result
	router.RegisterCallback(cmdSubmit, s.submit)

	// option updates
	router.RegisterCallbackPrefix(cmdUpdate, s.update)
	router.RegisterHandler(cmdGetData, s.getData)

	// media slider
	router.RegisterCallbackPrefix(cmdUpdateSlide, s.updateSlide)
	router.RegisterCallback(cmdCloseSlider, s.cancelSlider)
	// FIXME implemect cancel buttons

	// selector for schedule modify
	datetime.Register(
		router.With(cmdSelectMedia),
		scheduleAdd,
		session,
		s.canceledDateTimeSelector,
		onError,
		s.msgIdStorage,
		s.mediaSelected,
	)

	// null handler to answer callbacks for empty buttons
	router.RegisterCallback(cmdNoOp, s.nullHandler)
}

func (s *Search) init(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID

	if !s.auth.IsKnown(chatId) {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrUnknown,
		}); err != nil {
			s.onError(err)
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
		s.onError(err)
	}

	s.msgIdStorage.Set(chatId, msg.ID)
}

// submit gets text to search,
// search option and returns
// slider for found media.
func (s *Search) submit(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	opt := s.searchStorage.Get(chatId)

	tags := make([]string, 0)
	tags = append(tags, opt.format.String())
	tags = append(tags, opt.playlists...)
	tags = append(tags, opt.genres...)
	tags = append(tags, opt.languages...)
	tags = append(tags, opt.moods...)

	res, err := s.search.Search(chatId, localModels.MediaFilter{
		Name:       opt.nameAuthor,
		Author:     opt.nameAuthor,
		Tags:       tags,
		MaxRespLen: 20,
	})
	// TODO enhance errors
	if err != nil {
		s.onError(err)
	}

	if len(res) == 0 {
		if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatId,
			MessageID:   s.msgIdStorage.Get(chatId),
			Text:        ctr.LibSearchErrEmptyRes,
			ReplyMarkup: s.getSettingDataMarkup(),
		}); err != nil {
			s.onError(err)
		}
		return
	}

	s.mediaPage.Set(chatId, 1)
	s.mediaResults.Set(chatId, res)
	s.mediaSelected.Set(chatId, res[0])

	msg, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   s.msgIdStorage.Get(chatId),
		Text:        s.mediaRepr(res[0]),
		ParseMode:   models.ParseModeHTML,
		ReplyMarkup: s.mediaSliderMarkup(1, len(res)),
	})

	if err != nil {
		s.onError(err)
	}

	s.msgIdStorage.Set(chatId, msg.ID)
}

// TODO: update

func (s *Search) nullHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)
}

// filterRepr returns formated filter info.
func (s *Search) filterRepr(opt searchOption) string {
	var b strings.Builder

	b.WriteString("<b>Настройки поиска:</b>\n")
	if opt.nameAuthor != "" {
		b.WriteString(fmt.Sprintf("<b>Название/автор:</b> %s\n", opt.nameAuthor))
	}
	b.WriteString(fmt.Sprintf("<b>Формат:</b> %s\n", opt.format.String()))
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

// mediaRepr returns media formatted
// for telegram message.
func (s *Search) mediaRepr(media localModels.Media) string {
	var b strings.Builder

	b.WriteString("<b>Композиция</b>\n")
	b.WriteString(fmt.Sprintf("<b>Название:</b> %s\n", media.Name))
	b.WriteString(fmt.Sprintf("<b>Автор:</b> %s\n", media.Author))
	b.WriteString(fmt.Sprintf("<b>Длительность:</b> %s\n", media.Duration.Round(time.Second).String()))

	var podcasts, playlists, genres, languages, moods []string
	for _, tag := range media.Tags {
		switch tag.Type.Name {
		case "podcast":
			podcasts = append(podcasts, tag.Name)
		case "playlist":
			playlists = append(playlists, tag.Name)
		case "genre":
			genres = append(genres, tag.Name)
		case "language":
			languages = append(languages, tag.Name)
		case "mood":
			moods = append(moods, tag.Name)
		}
	}
	if len(podcasts) > 0 {
		b.WriteString(fmt.Sprintf("<b>Подкасты:</b> %s\n", strings.Join(podcasts, ", ")))
	}
	if len(playlists) > 0 {
		b.WriteString(fmt.Sprintf("<b>Плейлисты:</b> %s\n", strings.Join(playlists, ", ")))
	}
	if len(genres) > 0 {
		b.WriteString(fmt.Sprintf("<b>Жанры:</b> %s\n", strings.Join(genres, ", ")))
	}
	if len(languages) > 0 {
		b.WriteString(fmt.Sprintf("<b>Языки:</b> %s\n", strings.Join(languages, ", ")))
	}
	if len(moods) > 0 {
		b.WriteString(fmt.Sprintf("<b>Настроение:</b> %s\n", strings.Join(moods, ", ")))
	}

	return b.String()
}

func (s *Search) callbackAnswer(ctx context.Context, b *bot.Bot, callbackQuery *models.CallbackQuery) {
	ok, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
	})
	if err != nil {
		s.onError(err)
		return
	}
	if !ok {
		s.onError(fmt.Errorf("callback answer failed"))
	}
}

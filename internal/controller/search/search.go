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
	cmdCancel ctr.Command = "cancel"

	// options
	cmdNameAuthor    ctr.Command = "nameauthor"
	cmdFormat        ctr.Command = "format"
	cmdFormatAll     ctr.Command = "format-all"
	cmdFormatPodcast ctr.Command = "format-podcast"
	cmdFormatSong    ctr.Command = "format-song"
	cmdPlaylist      ctr.Command = "playlist"
	cmdGenre         ctr.Command = "genre"
	cmdLanguage      ctr.Command = "language"
	cmdMood          ctr.Command = "mood"
	cmdUpdateOption  ctr.Command = "update"

	// media slider
	cmdNextSlide   ctr.Command = "next-slide"
	cmdPrevSlide   ctr.Command = "prev-slide"
	cmdCloseSlider ctr.Command = "cancel"
	cmdSelectMedia ctr.Command = "select"

	// filler
	cmdNoOp ctr.Command = "no-op"
)

type Search struct {
	router   *ctr.Router
	search   LibrarySearch
	session  ctr.Session
	onCancel ctr.OnCancelHandler
	onError  bot.ErrorsHandler

	searchStorage       storage.Storage[searchOption]
	targetUpdateStorage storage.Storage[ctr.Command]
	mediaPage           storage.Storage[int]
	mediaResults        storage.Storage[[]localModels.Media]
	mediaSelected       storage.Storage[localModels.Media]
}

type LibrarySearch interface {
	Search(localModels.MediaFilter) ([]localModels.Media, error)
}

type ScheduleAdd interface {
	NewSegment(s localModels.Segment) error
}

type searchOption struct {
	nameAuthor string
	format     searchFormat
	playlists  []string
	genres     []string
	languages  []string
	moods      []string
}

type searchFormat string

const (
	formatAll     searchFormat = "все"
	formatSong    searchFormat = "песни"
	formatPodcast searchFormat = "подкасты"
)

var defaultOption = searchOption{format: formatAll}

func Register(
	router *ctr.Router,
	search LibrarySearch,
	scheduleAdd ScheduleAdd,
	session ctr.Session,
	onCancel ctr.OnCancelHandler,
	onError bot.ErrorsHandler,
) {
	s := &Search{
		router:   router,
		search:   search,
		session:  session,
		onCancel: onCancel,
		onError:  onError,

		searchStorage:       storage.Storage[searchOption]{},
		targetUpdateStorage: storage.Storage[ctr.Command]{},
		mediaPage:           storage.Storage[int]{},
		mediaResults:        storage.Storage[[]localModels.Media]{},
	}

	// main menu
	router.RegisterCallback(cmdBase, s.init)

	// call searcher and show result
	router.RegisterHandler(cmdSubmit, s.submit)

	// option updates
	router.RegisterCallback(cmdNameAuthor, s.nameAuthor)
	router.RegisterCallback(cmdFormat, s.format)
	router.RegisterCallback(cmdFormatAll, s.formatAll)
	router.RegisterCallback(cmdFormatPodcast, s.formatPodcast)
	router.RegisterCallback(cmdFormatSong, s.formatSong)
	router.RegisterCallback(cmdPlaylist, s.playlist)
	router.RegisterCallback(cmdGenre, s.genre)
	router.RegisterCallback(cmdLanguage, s.language)
	router.RegisterCallback(cmdMood, s.mood)

	// options update answer
	router.RegisterHandler(cmdUpdateOption, s.updateState)

	// media slider
	router.RegisterCallback(cmdNextSlide, s.nextMediaSlide)
	router.RegisterCallback(cmdPrevSlide, s.prevMediaSlide)
	router.RegisterCallback(cmdCloseSlider, s.updateState)

	// selector for schedule modify
	datetime.Register(
		router.With(cmdSelectMedia),
		scheduleAdd,
		session,
		s.canceledDateTimeSelector,
		onError,
		s.mediaSelected,
	)

	router.RegisterCallback(cmdCancel, s.cancelSearch)

	// null handler to answer callbacks for empty buttons
	router.RegisterCallback(cmdNoOp, s.nullHandler)
}

func (s *Search) init(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	s.searchStorage.Set(userId, defaultOption)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        ctr.LibSearchInit,
		ReplyMarkup: s.mainMenuMarkup(),
	}); err != nil {
		s.onError(err)
	}
}

// submit gets text to search,
// search option and returns
// slider for found media.
func (s *Search) submit(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID

	opt := s.searchStorage.Get(userId)

	tags := make([]string, 0)
	if opt.format != formatAll {
		tags = append(tags, string(opt.format))
	}
	tags = append(tags, opt.playlists...)
	tags = append(tags, opt.genres...)
	tags = append(tags, opt.languages...)
	tags = append(tags, opt.moods...)

	res, err := s.search.Search(localModels.MediaFilter{
		Name:       opt.nameAuthor,
		Author:     opt.nameAuthor,
		Tags:       tags,
		MaxRespLen: 20,
	})
	// TODO enhance errors
	if err != nil {
		s.onError(err)
	}

	s.mediaPage.Set(userId, 1)
	s.mediaResults.Set(userId, res)
	s.mediaSelected.Set(userId, res[0])

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        s.mediaRepr(res[0]),
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: s.mediaSliderMarkup(1, len(res)),
	}); err != nil {
		s.onError(err)
	}
}

func (s *Search) cancelSearch(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)

	s.onCancel(ctx, b, update.CallbackQuery.Message)
}

// TODO: update

func (s *Search) nullHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)
}

// filterRepr returns formated filter info.
func (s *Search) filterRepr(id int64) string {
	var b strings.Builder
	opt := s.searchStorage.Get(id)
	counter := 0

	b.WriteString("*Настройки поиска:*\n")
	if opt.nameAuthor != "" {
		b.WriteString(fmt.Sprintf("*Название/автор:* '%s'\n", opt.nameAuthor))
		counter++
	}
	if opt.format != formatAll {
		b.WriteString(fmt.Sprintf("*Формат:* %s\n", opt.format))
		counter++
	}
	if len(opt.playlists) > 0 {
		b.WriteString(fmt.Sprintf("*Плейлисты:* %s\n", strings.Join(opt.playlists, ", ")))
		counter++
	}
	if len(opt.genres) > 0 {
		b.WriteString(fmt.Sprintf("*Жанры:* %s\n", strings.Join(opt.genres, ", ")))
		counter++
	}
	if len(opt.languages) > 0 {
		b.WriteString(fmt.Sprintf("*Языки:* %s\n", strings.Join(opt.languages, ", ")))
		counter++
	}
	if len(opt.moods) > 0 {
		b.WriteString(fmt.Sprintf("*Настроения:* %s", strings.Join(opt.moods, ", ")))
		counter++
	}

	if counter == 0 {
		b.WriteString("Пока пустенько...")
	}

	return b.String()
}

// mediaRepr returns media formatted
// for telegram message.
func (s *Search) mediaRepr(media localModels.Media) string {
	var b strings.Builder

	b.WriteString("*Композиция*")
	b.WriteString(fmt.Sprintf("*Название:* %s\n", media.Name))
	b.WriteString(fmt.Sprintf("*Автор:* %s\n", media.Author))
	b.WriteString(fmt.Sprintf("*Длительность:* %s\n", media.Duration.Round(time.Second).String()))

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
		b.WriteString(fmt.Sprintf("*Подкасты*: %s\n", strings.Join(podcasts, ", ")))
	}
	if len(playlists) > 0 {
		b.WriteString(fmt.Sprintf("*Плейлисты*: %s\n", strings.Join(playlists, ", ")))
	}
	if len(genres) > 0 {
		b.WriteString(fmt.Sprintf("*Жанры*: %s\n", strings.Join(genres, ", ")))
	}
	if len(languages) > 0 {
		b.WriteString(fmt.Sprintf("*Языки*: %s\n", strings.Join(languages, ", ")))
	}
	if len(moods) > 0 {
		b.WriteString(fmt.Sprintf("*Настроение*: %s\n", strings.Join(moods, ", ")))
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

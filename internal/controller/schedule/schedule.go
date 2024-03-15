package schedule

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/storage"
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
)

const (
	// slider
	cmdNewPage ctr.Command = "new-slide"
	cmdUpdate  ctr.Command = "update"

	// filter
	cmdNoOp ctr.Command = "no-op"

	// page size
	pageSize = 10
)

type schedule struct {
	router  *ctr.Router
	auth    Auth
	sch     Schedule
	session ctr.Session
	onError bot.ErrorsHandler

	scheduleStorage  storage.Storage[[]localModels.Segment]
	respPagesStorage storage.Storage[int]
	msgIdStorage     storage.Storage[int]
}

type Auth interface {
	IsKnown(ctx context.Context, id int64) bool
}

type Schedule interface {
	Schedule(ctx context.Context, id int64) ([]localModels.Segment, error)
}

func Register(
	router *ctr.Router,
	auth Auth,
	sch Schedule,
	session ctr.Session,
	onError bot.ErrorsHandler,
) {
	s := &schedule{
		router:  router,
		auth:    auth,
		sch:     sch,
		session: session,
		onError: onError,

		scheduleStorage:  storage.New[[]localModels.Segment](),
		respPagesStorage: storage.New[int](),
		msgIdStorage:     storage.New[int](),
	}

	router.RegisterCommand(s.init)

	// slider
	router.RegisterCallbackPrefix(cmdNewPage, s.newPage)
	router.RegisterCallback(cmdUpdate, s.update)

	// filler
	router.RegisterCallback(cmdNoOp, s.nullHandler)
}

func (s *schedule) init(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "schedule.init"

	chatId := update.Message.Chat.ID

	if !s.auth.IsKnown(ctx, chatId) {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrUnknown,
		}); err != nil {
			s.onError(fmt.Errorf("%s: %w", op, err))
		}
		return
	}

	res, err := s.sch.Schedule(ctx, chatId)
	if err != nil {
		// handle errors
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			s.onError(fmt.Errorf("%s: %w", op, err))
		}
		return
	}

	pages := len(res) / pageSize
	if len(res)%pageSize != 0 {
		pages++
	}

	s.respPagesStorage.Set(chatId, pages)
	s.scheduleStorage.Set(chatId, res)

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        s.scheduleFormat(res, 1),
		ReplyMarkup: s.mainMenuMarkup(1, pages),
		ParseMode:   models.ParseModeHTML,
	})
	if err != nil {
		s.onError(fmt.Errorf("%s: %w", op, err))
	}

	s.msgIdStorage.Set(chatId, msg.ID)
}

func (s *schedule) update(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "schedule.update"

	s.callbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	res, err := s.sch.Schedule(ctx, chatId)
	if err != nil {
		// handle errors
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			s.onError(fmt.Errorf("%s: %w", op, err))
		}
		return
	}

	pages := len(res) / pageSize
	if len(res)%pageSize != 0 {
		pages++
	}

	s.respPagesStorage.Set(chatId, pages)
	s.scheduleStorage.Set(chatId, res)

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   s.msgIdStorage.Get(chatId),
		Text:        s.scheduleFormat(res, 1),
		ReplyMarkup: s.mainMenuMarkup(1, pages),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		s.onError(fmt.Errorf("%s: %w", op, err))
	}
}

func (s *schedule) newPage(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "schedule.newPage"

	s.callbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	id, err := strconv.Atoi(s.router.GetState(update.CallbackQuery.Data))
	if err != nil {
		if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    chatId,
			MessageID: s.msgIdStorage.Get(chatId),
			Text:      ctr.ErrorMessage,
		}); err != nil {
			s.onError(fmt.Errorf("%s: %w", op, err))
		}
		return
	}

	res := s.scheduleStorage.Get(chatId)
	pages := s.respPagesStorage.Get(chatId)

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   s.msgIdStorage.Get(chatId),
		Text:        s.scheduleFormat(res, id),
		ReplyMarkup: s.mainMenuMarkup(id, pages),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		s.onError(fmt.Errorf("%s: %w", op, err))
	}
}

func (s *schedule) nullHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	s.callbackAnswer(ctx, b, update.CallbackQuery)
}

func (s *schedule) callbackAnswer(ctx context.Context, b *bot.Bot, callbackQuery *models.CallbackQuery) {
	const op = "schedule.callbackAnswer"

	ok, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
	})
	if err != nil {
		s.onError(fmt.Errorf("%s: %w", op, err))
		return
	}
	if !ok {
		s.onError(fmt.Errorf("callback answer failed"))
	}
}

func (s *schedule) scheduleFormat(sch []localModels.Segment, page int) string {
	var b strings.Builder

	startId := (page - 1) * pageSize
	stopId := min(len(sch), page*pageSize)

	// TODO: highlight protected segments.
	for _, s := range sch[startId:stopId] {
		b.WriteString(fmt.Sprintf(
			"[%s-%s]\n%s \u2014 %s\n",
			s.Start.Format("2006-01-02"),
			s.Start.Add(s.StopCut-s.BeginCut).Format("2006-01-02"),
			s.Media.Name,
			s.Media.Author,
		))
	}

	return b.String()
}

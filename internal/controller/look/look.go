package look

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
	cmdBase ctr.Command = ""
	cmdBack ctr.Command = "back"

	// slider
	cmdNewPage ctr.Command = "new-slide"
	cmdUpdate  ctr.Command = "update"

	// filter
	cmdNoOp ctr.Command = "no-op"

	// page size
	pageSize = 20
)

type look struct {
	router   *ctr.Router
	sch      Schedule
	session  ctr.Session
	onCancel ctr.OnCancelHandler
	onError  bot.ErrorsHandler

	scheduleStorage  storage.Storage[[]localModels.Segment]
	respPagesStorage storage.Storage[int]
}

type Schedule interface {
	Schedule(ctx context.Context) ([]localModels.Segment, error)
}

func Register(
	router *ctr.Router,
	sch Schedule,
	session ctr.Session,
	onCancel ctr.OnCancelHandler,
	onError bot.ErrorsHandler,
) {
	l := &look{
		router:   router,
		sch:      sch,
		session:  session,
		onCancel: onCancel,
		onError:  onError,

		scheduleStorage: storage.Storage[[]localModels.Segment]{},
	}

	router.RegisterCallback(cmdBase, l.init)

	// slider
	router.RegisterCallbackPrefix(cmdNewPage, l.newPage)
	router.RegisterCallback(cmdUpdate, l.init)
	router.RegisterCallbackPrefix(cmdBack, l.cancelPage)

	// filler
	router.RegisterCallback(cmdNoOp, l.nullHandler)
}

func (l *look) init(ctx context.Context, b *bot.Bot, update *models.Update) {
	l.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	res, err := l.sch.Schedule(ctx)
	if err != nil {
		// handle errors
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			l.onError(err)
		}
		return
	}

	pages := len(res) / pageSize
	if len(res)%pageSize != 0 {
		pages++
	}

	l.respPagesStorage.Set(userId, pages)
	l.scheduleStorage.Set(userId, res)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        l.scheduleFormat(res, 1),
		ReplyMarkup: l.mainMenuMarkup(1, pages),
		ParseMode:   models.ParseModeMarkdown,
	}); err != nil {
		l.onError(err)
	}
}

func (l *look) newPage(ctx context.Context, b *bot.Bot, update *models.Update) {
	l.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	id, err := strconv.Atoi(l.router.GetState(update.CallbackQuery.Data))
	if err != nil {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			l.onError(err)
		}
		return
	}

	res := l.scheduleStorage.Get(userId)
	pages := l.respPagesStorage.Get(userId)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        l.scheduleFormat(res, id),
		ReplyMarkup: l.mainMenuMarkup(id, pages),
		ParseMode:   models.ParseModeMarkdown,
	}); err != nil {
		l.onError(err)
	}
}

func (l *look) cancelPage(ctx context.Context, b *bot.Bot, update *models.Update) {
	l.callbackAnswer(ctx, b, update.CallbackQuery)

	l.onCancel(ctx, b, update.CallbackQuery.Message)
}

func (l *look) nullHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	l.callbackAnswer(ctx, b, update.CallbackQuery)
}

func (l *look) callbackAnswer(ctx context.Context, b *bot.Bot, callbackQuery *models.CallbackQuery) {
	ok, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
	})
	if err != nil {
		l.onError(err)
		return
	}
	if !ok {
		l.onError(fmt.Errorf("callback answer failed"))
	}
}

func (l *look) scheduleFormat(sch []localModels.Segment, page int) string {
	var b strings.Builder

	// TODO: highlight protected segments.
	for _, s := range sch[(page-1)*pageSize : page*pageSize] {
		b.WriteString(fmt.Sprintf(
			"[%s-%s]\n%s \u2014 %s",
			s.Start.Format("2006-01-02"),
			s.Start.Add(s.StopCut-s.BeginCut).Format("2006-01-02"),
			s.Media.Name,
			s.Media.Author,
		))
	}

	return b.String()
}

package datetime

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/datepicker"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/storage"
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
)

const (
	cmdBase ctr.Command = ""

	cmdDate   ctr.Command = "date"
	cmdSubmit ctr.Command = "submit"

	cmdCancelTime ctr.Command = "cancel"
)

type picker struct {
	router       *ctr.Router
	schedule     ScheduleAdd
	session      ctr.Session
	onCancel     ctr.OnCancelHandler
	onError      bot.ErrorsHandler
	mediaStorage storage.Storage[localModels.Media]

	dateStorage        storage.Storage[time.Time]
	pickMessageStorage storage.Storage[int]
}

type ScheduleAdd interface {
	NewSegment(s localModels.Segment) error
}

func Register(
	router *ctr.Router,
	schedule ScheduleAdd,
	session ctr.Session,
	onCancel ctr.OnCancelHandler,
	onError bot.ErrorsHandler,
	mediaStorage storage.Storage[localModels.Media],
) {
	p := &picker{
		router:   router,
		schedule: schedule,
		session:  session,
		onCancel: onCancel,
		onError:  onError,

		dateStorage:        storage.Storage[time.Time]{},
		pickMessageStorage: storage.Storage[int]{},
		mediaStorage:       mediaStorage,
	}

	router.RegisterCallback(cmdBase, p.init)
	router.RegisterCallbackPrefix(cmdSubmit, p.submitDateTime)
	router.RegisterCallback(cmdCancelTime, p.cancelTimePicker)
}

func (p *picker) init(ctx context.Context, b *bot.Bot, update *models.Update) {
	p.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   ctr.LibSearchInit,
		ReplyMarkup: datepicker.New(
			b, p.catchDatePicker,
			datepicker.CurrentDate(time.Now()),
			datepicker.From(time.Now()),
			datepicker.WithPrefix(p.router.Path(cmdDate)),
			datepicker.OnCancel(datepicker.OnCancelHandler(p.onCancel)),
			datepicker.OnError(datepicker.OnErrorHandler(p.onError)),
		),
	})
	if err != nil {
		p.onError(err)
		return
	}
	p.pickMessageStorage.Set(userId, msg.ID)
}

func (p *picker) catchDatePicker(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, date time.Time) {
	userId := mes.Message.From.ID
	chatId := mes.Message.Chat.ID

	p.dateStorage.Set(userId, date)

	if _, err := b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:      chatId,
		MessageID:   p.pickMessageStorage.Get(userId),
		ReplyMarkup: p.timePickMarkup(),
	}); err != nil {
		p.onError(err)
	}
}

func (p *picker) submitDateTime(ctx context.Context, b *bot.Bot, update *models.Update) {
	p.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	y, m, d := p.dateStorage.Get(userId).Date()
	timeStr := p.router.GetState(update.CallbackQuery.Data)

	H, M, _ := strings.Cut(timeStr, ":")
	hour, _ := strconv.Atoi(H)
	minute, _ := strconv.Atoi(M)

	date := time.Date(y, m, d, hour, minute, 0, 0, time.Local)

	media := p.mediaStorage.Get(userId)
	segm := localModels.Segment{
		Media:     media,
		Start:     date,
		BeginCut:  0,
		StopCut:   media.Duration,
		Protected: true,
	}
	if err := p.schedule.NewSegment(segm); err != nil {
		p.onError(err)
		// TODO: handle many errors
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   p.successMsg(date, date.Add(media.Duration)),
	}); err != nil {
		p.onError(err)
	}
}

func (p *picker) cancelTimePicker(ctx context.Context, b *bot.Bot, update *models.Update) {
	p.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	if _, err := b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:    chatId,
		MessageID: p.pickMessageStorage.Get(userId),
		ReplyMarkup: datepicker.New(
			b, p.catchDatePicker,
			datepicker.CurrentDate(time.Now()),
			datepicker.From(time.Now()),
			datepicker.WithPrefix(p.router.Path(cmdDate)),
			datepicker.OnCancel(datepicker.OnCancelHandler(p.onCancel)),
			datepicker.OnError(datepicker.OnErrorHandler(p.onError)),
		),
	}); err != nil {
		p.onError(err)
	}
}

func (p *picker) callbackAnswer(ctx context.Context, b *bot.Bot, callbackQuery *models.CallbackQuery) {
	ok, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
	})
	if err != nil {
		p.onError(err)
		return
	}
	if !ok {
		p.onError(fmt.Errorf("callback answer failed"))
	}
}

func (p *picker) successMsg(start, stop time.Time) string {
	return fmt.Sprintf("Добавлено в расписание с %s по %s.", start.Format("15:04:05"), stop.Format("15:04:05"))
}
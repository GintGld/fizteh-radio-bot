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

	dateStorage  storage.Storage[time.Time]
	msgIdStorage storage.Storage[int]
}

type ScheduleAdd interface {
	NewSegment(ctx context.Context, id int64, s localModels.Segment) error
}

func Register(
	router *ctr.Router,
	schedule ScheduleAdd,
	session ctr.Session,
	onCancel ctr.OnCancelHandler,
	onError bot.ErrorsHandler,
	msgIdStorage storage.Storage[int],
	mediaStorage storage.Storage[localModels.Media],
) {
	p := &picker{
		router:   router,
		schedule: schedule,
		session:  session,
		onCancel: onCancel,
		onError:  onError,

		dateStorage:  storage.New[time.Time](),
		msgIdStorage: msgIdStorage,
		mediaStorage: mediaStorage,
	}

	router.RegisterCallback(cmdBase, p.init)
	router.RegisterCallbackPrefix(cmdSubmit, p.submitDateTime)
	router.RegisterCallback(cmdCancelTime, p.cancelTimePicker)
}

func (p *picker) init(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "picker.init"

	p.callbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	msg, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    chatId,
		MessageID: p.msgIdStorage.Get(chatId),
		Text:      ctr.LibSearchInit,
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
		p.onError(fmt.Errorf("%s: %w", op, err))
		return
	}
	p.msgIdStorage.Set(chatId, msg.ID)
}

func (p *picker) catchDatePicker(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage, date time.Time) {
	const op = "picker.catchDatePicker"

	chatId := mes.Message.Chat.ID

	p.dateStorage.Set(chatId, date)

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		Text:        ctr.LibSearchPickSelecting,
		ChatID:      chatId,
		ReplyMarkup: p.timePickMarkup(),
	})
	if err != nil {
		p.onError(fmt.Errorf("%s: %w", op, err))
	}

	p.msgIdStorage.Set(chatId, msg.ID)
}

func (p *picker) submitDateTime(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "picker.submitDateTime"

	p.callbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	y, m, d := p.dateStorage.Get(chatId).Date()
	timeStr := p.router.GetState(update.CallbackQuery.Data)

	H, M, _ := strings.Cut(timeStr, ":")
	hour, _ := strconv.Atoi(H)
	minute, _ := strconv.Atoi(M)

	date := time.Date(y, m, d, hour, minute, 0, 0, time.Local)

	media := p.mediaStorage.Get(chatId)
	segm := localModels.Segment{
		Media:     media,
		Start:     date,
		BeginCut:  0,
		StopCut:   media.Duration,
		Protected: true,
	}
	if err := p.schedule.NewSegment(ctx, chatId, segm); err != nil {
		p.onError(fmt.Errorf("%s: %w", op, err))
		// TODO: handle many errors
	}

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    chatId,
		MessageID: p.msgIdStorage.Get(chatId),
		Text:      p.successMsg(date, date.Add(media.Duration)),
	}); err != nil {
		p.onError(fmt.Errorf("%s: %w", op, err))
	}
}

func (p *picker) cancelTimePicker(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "picker.cancelTimePicker"

	p.callbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	if _, err := b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:    chatId,
		MessageID: p.msgIdStorage.Get(chatId),
		ReplyMarkup: datepicker.New(
			b, p.catchDatePicker,
			datepicker.CurrentDate(time.Now()),
			datepicker.From(time.Now()),
			datepicker.WithPrefix(p.router.Path(cmdDate)),
			datepicker.OnCancel(datepicker.OnCancelHandler(p.onCancel)),
			datepicker.OnError(datepicker.OnErrorHandler(p.onError)),
		),
	}); err != nil {
		p.onError(fmt.Errorf("%s: %w", op, err))
	}
}

func (p *picker) callbackAnswer(ctx context.Context, b *bot.Bot, callbackQuery *models.CallbackQuery) {
	const op = "picker.callbackAnswer"

	ok, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
	})
	if err != nil {
		p.onError(fmt.Errorf("%s: %w", op, err))
		return
	}
	if !ok {
		p.onError(fmt.Errorf("callback answer failed"))
	}
}

func (p *picker) successMsg(start, stop time.Time) string {
	return fmt.Sprintf("Добавлено в расписание с %s по %s.", start.Format("15:04:05"), stop.Format("15:04:05"))
}

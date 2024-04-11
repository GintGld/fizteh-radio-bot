package live

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/storage"
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
)

const (
	cmdStart      ctr.Command = "start"
	cmdStop       ctr.Command = "stop"
	cmdStopSubmit ctr.Command = "submit"
	cmdStopReject ctr.Command = "reject"
	cmdGetName    ctr.Command = "name"
)

type live struct {
	ctr.CallbackAnswerer

	router  *ctr.Router
	auth    Auth
	live    LiveSrv
	session ctr.Session
	onError bot.ErrorsHandler

	msgIdStorage storage.Storage[int]
}

type Auth interface {
	IsKnown(ctx context.Context, id int64) bool
}

type LiveSrv interface {
	StartLive(ctx context.Context, id int64, live localModels.Live) error
	StopLive(ctx context.Context, id int64) error
	LiveInfo(ctx context.Context, id int64) (localModels.Live, error)
}

func Register(
	router *ctr.Router,
	auth Auth,
	liveSrv LiveSrv,
	session ctr.Session,
	onError bot.ErrorsHandler,
) {
	l := &live{
		router:  router,
		auth:    auth,
		live:    liveSrv,
		session: session,
		onError: onError,

		msgIdStorage: storage.New[int](),
	}

	router.RegisterCommand(l.init)

	router.RegisterCallback(cmdStart, l.start)
	router.RegisterHandler(cmdGetName, l.getName)

	router.RegisterCallback(cmdStop, l.stop)
	router.RegisterCallback(cmdStopSubmit, l.submitStop)
	router.RegisterCallback(cmdStopReject, l.rejectStop)
}

func (l *live) init(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "live.init"

	chatId := update.Message.Chat.ID

	if !l.auth.IsKnown(ctx, chatId) {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrUnknown,
		}); err != nil {
			l.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	var (
		msgText string
		markup  models.InlineKeyboardMarkup
	)

	live, err := l.live.LiveInfo(ctx, chatId)
	if err != nil {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			l.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}
	if live.ID == 0 {
		msgText = ctr.LiveNotPlaying
		markup = l.StartMarkup()
	} else {
		msgText = live.String()
		markup = l.StopMarkup()
	}

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        msgText,
		ReplyMarkup: markup,
		ParseMode:   models.ParseModeHTML,
	})
	if err != nil {
		l.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}

	l.msgIdStorage.Set(chatId, msg.ID)
}

func (l *live) start(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "live.start"

	l.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	l.session.Redirect(chatId, l.router.Path(cmdGetName))

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		MessageID: l.msgIdStorage.Get(chatId),
		ChatID:    chatId,
		Text:      ctr.LiveAskName,
	})
	if err != nil {
		l.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (l *live) getName(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "live.getName"

	chatId := update.Message.Chat.ID
	l.session.Redirect(chatId, ctr.NullStatus)

	name := update.Message.Text
	if name == "" {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.LiveNameEmpty,
		}); err != nil {
			l.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	if err := l.live.StartLive(ctx, chatId, localModels.Live{Name: name}); err != nil {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			l.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    chatId,
		MessageID: update.Message.ID,
	})

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		MessageID: l.msgIdStorage.Get(chatId),
		ChatID:    chatId,
		Text:      ctr.LiveStarted,
	})
	if err != nil {
		l.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (l *live) stop(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "live.stop"

	l.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		MessageID:   l.msgIdStorage.Get(chatId),
		ChatID:      chatId,
		Text:        ctr.LiveSubmitStop,
		ReplyMarkup: l.SubmitStopMarkup(),
	})
	if err != nil {
		l.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (l *live) submitStop(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "live.submitStop"

	l.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	if err := l.live.StopLive(ctx, chatId); err != nil {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			l.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		MessageID: l.msgIdStorage.Get(chatId),
		ChatID:    chatId,
		Text:      ctr.LiveStopped,
	})
	if err != nil {
		l.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (l *live) rejectStop(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "live.rejectStop"

	l.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	var (
		msgText string
		markup  models.InlineKeyboardMarkup
	)

	live, err := l.live.LiveInfo(ctx, chatId)
	if err != nil {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			l.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}
	if live.ID == 0 {
		msgText = ctr.LiveNotPlaying
		markup = l.StartMarkup()
	} else {
		msgText = live.String()
		markup = l.StopMarkup()
	}

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   l.msgIdStorage.Get(chatId),
		Text:        msgText,
		ReplyMarkup: markup,
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		l.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

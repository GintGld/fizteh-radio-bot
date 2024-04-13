package stat

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/storage"
)

type stat struct {
	ctr.CallbackAnswerer

	router  *ctr.Router
	auth    Auth
	stat    Stat
	onError bot.ErrorsHandler

	msgIdStorage storage.Storage[int]
}

type Auth interface {
	IsKnown(ctx context.Context, id int64) bool
}

type Stat interface {
	ListenersNumber(ctx context.Context, id int64) (int64, error)
}

func Register(
	router *ctr.Router,
	auth Auth,
	statSrv Stat,
	onError bot.ErrorsHandler,
) {
	s := &stat{
		router:  router,
		auth:    auth,
		stat:    statSrv,
		onError: onError,

		msgIdStorage: storage.New[int](),
	}

	router.RegisterCommand(s.init)
}

func (s *stat) init(ctx context.Context, b *bot.Bot, update *models.Update) {
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

	N, err := s.stat.ListenersNumber(ctx, chatId)
	if err != nil {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatId,
		Text:      s.formatListenersNumber(N),
		ParseMode: models.ParseModeHTML,
	})

	if err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}

	s.msgIdStorage.Set(chatId, msg.ID)
}

func (s *stat) formatListenersNumber(N int64) string {
	return fmt.Sprintf("<b>Сейчас слушают</b>: %d", N)
}

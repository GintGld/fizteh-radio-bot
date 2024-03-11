package schedule

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/GintGld/fizteh-radio-bot/internal/controller/autodj"
	"github.com/GintGld/fizteh-radio-bot/internal/controller/look"
)

const (
	cmdLook = "look"
	cmdDj   = "dj"
)

type schedule struct {
	router  *ctr.Router
	auth    Auth
	session ctr.Session
	onError bot.ErrorsHandler
}

type Auth interface {
	IsKnown(id int64) bool
	Login(login, pass string) error
}

func Register(
	router *ctr.Router,
	auth Auth,
	sch look.Schedule,
	dj autodj.AutoDJ,
	session ctr.Session,
	onError bot.ErrorsHandler,
) {
	s := &schedule{
		router:  router,
		auth:    auth,
		session: session,
		onError: onError,
	}

	router.RegisterCommand(s.init)

	look.Register(
		router.With(cmdLook),
		sch,
		session,
		s.cancelSubmodule,
		onError,
	)

	autodj.Register(
		router.With(cmdDj),
		dj,
		session,
		s.cancelSubmodule,
		onError,
	)

}

func (s *schedule) init(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID

	if !s.auth.IsKnown(userId) {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrUnknown,
		}); err != nil {
			s.onError(err)
		}
		return
	}

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        ctr.SchMainMenuMessage,
		ReplyMarkup: s.mainMenuMarkup(),
	}); err != nil {
		s.onError(err)
	}
}

func (s *schedule) cancelSubmodule(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage) {
	chatId := mes.Message.Chat.ID

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        ctr.SchMainMenuMessage,
		ReplyMarkup: s.mainMenuMarkup(),
	}); err != nil {
		s.onError(err)
	}
}

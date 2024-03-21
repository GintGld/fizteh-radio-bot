package start

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/storage"
)

const (
	cmdLogin = "login"
	cmdPass  = "pass"
)

// TODO delete messages with login and password after authorization.

type start struct {
	router  *ctr.Router
	auth    Auth
	session ctr.Session
	onError bot.ErrorsHandler

	loginStorage storage.Storage[string]
	msgToDel     storage.Storage[[]int]
}

type Auth interface {
	IsKnown(ctx context.Context, id int64) bool
	Login(ctx context.Context, id int64, login, pass string) error
}

func Register(
	router *ctr.Router,
	auth Auth,
	session ctr.Session,
	onError bot.ErrorsHandler,
) {
	app := &start{
		router:  router,
		auth:    auth,
		session: session,
		onError: onError,

		loginStorage: storage.New[string](),
		msgToDel:     storage.New[[]int](),
	}

	router.RegisterCommand(app.init)
	router.RegisterHandler(cmdLogin, app.login)
	router.RegisterHandler(cmdPass, app.pass)
}

// First what user will see.
//
// May redirect to authorization
// if user is unknown.
func (s *start) init(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "start.init"

	chatId := update.Message.Chat.ID

	// Check if user is known or not
	if s.auth.IsKnown(ctx, update.Message.From.ID) {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   fmt.Sprintf(ctr.AuthorizedMessage, update.Message.From.FirstName),
		}); err != nil {
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
	} else {
		s.session.Redirect(chatId, s.router.Path(cmdLogin))
		msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.HelloMessage,
		})
		if err != nil {
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
			return
		}
		s.msgToDel.Set(chatId, []int{msg.ID})
	}
}

// Get login, validate it, ask for a password.
func (s *start) login(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "start.login"

	chatId := update.Message.Chat.ID

	login := update.Message.Text
	if login == "" {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrEmptyLogin,
		}); err != nil {
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	s.loginStorage.Set(chatId, login)
	s.session.Redirect(chatId, s.router.Path(cmdPass))

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   ctr.GotLoginAskPass,
	})
	if err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		return
	}
	msgs := s.msgToDel.Get(chatId)
	msgs = append(msgs, msg.ID)
	msgs = append(msgs, update.Message.ID)
	s.msgToDel.Set(chatId, msgs)
}

// Get pass, validate it
func (s *start) pass(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "start.pass"

	chatId := update.Message.Chat.ID

	pass := update.Message.Text
	if pass == "" {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrEmptyPass,
		}); err != nil {
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	login := s.loginStorage.Get(chatId)

	if err := s.auth.Login(ctx, chatId, login, pass); err != nil {
		// TODO: error statuses
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrAuthorizedMessage,
		}); err != nil {
			s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		s.session.Redirect(chatId, s.router.Path(cmdLogin))
		return
	}
	msgs := s.msgToDel.Get(chatId)
	msgs = append(msgs, update.Message.ID)

	if _, err := b.DeleteMessages(ctx, &bot.DeleteMessagesParams{
		ChatID:     chatId,
		MessageIDs: msgs,
	}); err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		return
	}

	s.msgToDel.Del(chatId)
	s.loginStorage.Del(chatId)
	s.session.Redirect(chatId, ctr.NullStatus)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   fmt.Sprintf(ctr.WelcomeMessage, update.Message.From.FirstName),
	}); err != nil {
		s.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

package start

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
)

const (
	keyLogin = "login"

	cmdLogin = "/login"
	cmdPass  = "/pass"
)

type Start struct {
	router  *ctr.Router
	log     *slog.Logger
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
	log *slog.Logger,
	auth Auth,
	session ctr.Session,
	onError bot.ErrorsHandler,
) {
	app := &Start{
		router:  router,
		log:     log,
		auth:    auth,
		session: session,
		onError: onError,
	}

	router.RegisterCommand(app.init)
	router.RegisterHandler(cmdLogin, app.login)
	router.RegisterHandler(cmdPass, app.pass)
}

// First what user will see.
//
// May redirect to authorization
// if user is unknown.
func (s *Start) init(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID

	// Check if user is known or not
	if s.auth.IsKnown(update.Message.From.ID) {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   fmt.Sprintf(ctr.AuthorizedMessage, update.Message.From.FirstName),
		}); err != nil {
			s.onError(err)
		}
	} else {
		s.session.Redirect(userId, s.router.FullPath(cmdLogin))
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.HelloMessage,
		}); err != nil {
			s.onError(err)
		}
	}
}

// Get login, validate it, ask for a password.
func (s *Start) login(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID

	login := update.Message.Text
	if login == "" {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrEmptyLogin,
		}); err != nil {
			s.onError(err)
		}
		return
	}

	s.session.Set(userId, keyLogin, login)
	s.session.Redirect(userId, s.router.FullPath(cmdPass))

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   ctr.GotLoginAskPass,
	}); err != nil {
		s.onError(err)
	}
}

// Get pass, validate it
func (s *Start) pass(ctx context.Context, b *bot.Bot, update *models.Update) {
	userId := update.Message.From.ID
	chatId := update.Message.Chat.ID

	pass := update.Message.Text
	if pass == "" {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrEmptyPass,
		}); err != nil {
			s.onError(err)
		}
		return
	}

	login := s.session.Get(userId, keyLogin)

	if err := s.auth.Login(login, pass); err != nil {
		// TODO: error statuses
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrAuthorizedMessage,
		}); err != nil {
			s.onError(err)
		}
		return
	}

	s.session.Del(userId, login)
	s.session.Redirect(userId, ctr.NullStatus)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   fmt.Sprintf(ctr.WelcomeMessage, update.Message.From.FirstName),
	}); err != nil {
		s.onError(err)
	}
}

package auth

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/GintGld/fizteh-radio-bot/internal/lib/logger/sl"
	"github.com/GintGld/fizteh-radio-bot/internal/models"
	"github.com/GintGld/fizteh-radio-bot/internal/service"

	"github.com/golang-jwt/jwt/v5"
)

// TODO dump user info and recover it

type auth struct {
	log        *slog.Logger
	authClient AuthClient
	users      map[int64]models.User
	userMutex  map[int64]*sync.Mutex
	updCtx     context.Context
	cancel     context.CancelFunc
}

type AuthClient interface {
	GetToken(ctx context.Context, user models.User) (jwt.Token, error)
}

func New(
	log *slog.Logger,
	authCLient AuthClient,
) *auth {
	ctx, cancel := context.WithCancel(context.Background())

	return &auth{
		log:        log,
		authClient: authCLient,
		users:      make(map[int64]models.User),
		userMutex:  make(map[int64]*sync.Mutex),
		updCtx:     ctx,
		cancel:     cancel,
	}
}

func (a *auth) IsKnown(_ context.Context, id int64) bool {
	_, ok := a.users[id]

	return ok
}

// Login logins user and
// setup user token update.
func (a *auth) Login(ctx context.Context, id int64, login, pass string) error {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
	)

	token, err := a.authClient.GetToken(ctx, models.User{
		Login: login,
		Pass:  pass,
	})
	if err != nil {
		// TODO handle errors
		log.Error(
			"failed to get token",
			slog.Int64("id", id),
			slog.String("login", login),
			sl.Err(err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	a.users[id] = models.User{
		Login: login,
		Pass:  pass,
		Token: token,
	}
	a.userMutex[id] = &sync.Mutex{}

	a.updateToken(a.updCtx, id)

	panic("not implemented")
}

// Token returns user's token
// if user does not exists
// returns service.ErrUserNotFound error.
func (a *auth) Token(_ context.Context, id int64) (jwt.Token, error) {
	a.userMutex[id].Lock()
	defer a.userMutex[id].Unlock()

	if user, ok := a.users[id]; ok {
		return user.Token, nil
	} else {
		return jwt.Token{}, service.ErrUserNotFound
	}
}

func (a *auth) Dump() {
	// TODO save user info to a file

	a.cancel()
}

// updateToken updates token for user.
func (a *auth) updateToken(ctx context.Context, id int64) error {
	const op = "auth.updateToken"
	const timeUntilTokenExpires = 5 * time.Second

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	a.userMutex[id].Lock()
	defer a.userMutex[id].Unlock()

	log := a.log.With(
		slog.String("op", op),
	)

	user := a.users[id]

	token, err := a.authClient.GetToken(ctx, user)
	if err != nil {
		// TODO handle errors
		log.Error(
			"failed to get token",
			slog.Int64("id", id),
			slog.String("login", user.Login),
			sl.Err(err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		log.Error(
			"failed to get expiration date",
			slog.Int64("id", id),
			sl.Err(err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	a.users[id] = models.User{
		Login: user.Login,
		Pass:  user.Pass,
		Token: token,
	}

	time.AfterFunc(time.Until(exp.Time)-timeUntilTokenExpires, func() {
		a.updateToken(ctx, id)
	})

	return nil
}

package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
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
	cacheFile  string

	users     map[int64]models.User
	mapMutex  *sync.Mutex
	userMutex map[int64]*sync.Mutex
	updCtx    context.Context
	cancel    context.CancelFunc
}

type AuthClient interface {
	GetToken(ctx context.Context, user models.User) (jwt.Token, error)
}

func New(
	log *slog.Logger,
	authCLient AuthClient,
	cacheFile string,
) *auth {
	ctx, cancel := context.WithCancel(context.Background())

	a := &auth{
		log:        log,
		authClient: authCLient,
		cacheFile:  cacheFile,
		users:      make(map[int64]models.User),
		mapMutex:   &sync.Mutex{},
		userMutex:  make(map[int64]*sync.Mutex),
		updCtx:     ctx,
		cancel:     cancel,
	}

	if err := a.recoverUsers(context.Background()); err != nil {
		log.Error(
			"failed to recover users",
			slog.String("op", "auth.New"),
			sl.Err(err),
		)
	}

	return a
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

	if err := a.Dump(); err != nil {
		log.Error(
			"failed to dump new users info",
			sl.Err(err),
		)
	}

	return nil
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

func (a *auth) recoverUsers(ctx context.Context) error {
	const op = "auth.recoverUsers"

	log := a.log.With(
		slog.String("op", op),
		slog.String("file", a.cacheFile),
	)

	file, err := os.Open(a.cacheFile)
	if err != nil {
		log.Error("failed to open cache file", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := io.ReadAll(file)
	if err != nil {
		log.Error("failed to open cache file", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := json.Unmarshal(res, &a.users); err != nil {
		log.Error("failed to unmarshal users info", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	for id, user := range a.users {
		if err := a.Login(ctx, id, user.Login, user.Pass); err != nil {
			log.Error("failed to login user", slog.Int64("id", id), sl.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (a *auth) Dump() error {
	const op = "auth.Dump"

	log := a.log.With(
		slog.String("op", op),
		slog.String("file", a.cacheFile),
	)

	a.mapMutex.Lock()
	defer a.mapMutex.Unlock()

	bytes, err := json.Marshal(a.users)
	if err != nil {
		log.Error("failed to marshal users", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	file, err := os.OpenFile(a.cacheFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("failed to open file", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := file.Write(bytes); err != nil {
		log.Error("failed to write to a file", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
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

package stat

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/GintGld/fizteh-radio-bot/internal/lib/logger/sl"
	"github.com/golang-jwt/jwt/v5"
)

type stat struct {
	log        *slog.Logger
	auth       Auth
	statClient StatClient
}

type Auth interface {
	Token(ctx context.Context, id int64) (jwt.Token, error)
}

type StatClient interface {
	ListenersNumber(ctx context.Context) (int64, error)
}

func New(
	log *slog.Logger,
	auth Auth,
	statClient StatClient,
) *stat {
	return &stat{
		log:        log,
		auth:       auth,
		statClient: statClient,
	}
}

func (s *stat) ListenersNumber(ctx context.Context, id int64) (int64, error) {
	const op = "stat.ListenersNumber"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("id", id),
	)

	N, err := s.statClient.ListenersNumber(ctx)
	if err != nil {
		log.Error("failed to get liesteners number", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return N, nil
}

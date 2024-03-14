package auth

import (
	"context"
	"log/slog"
)

type auth struct {
	log *slog.Logger
}

func New(
	log *slog.Logger,
) *auth {
	return &auth{
		log: log,
	}
}

func (a *auth) IsKnown(ctx context.Context, id int64) bool {
	panic("not implemented")
}

func (a *auth) Login(ctx context.Context, id int64, login, pass string) error {
	panic("not implemented")
}

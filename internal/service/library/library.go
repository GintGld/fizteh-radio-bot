package library

import (
	"context"
	"log/slog"

	"github.com/GintGld/fizteh-radio-bot/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type library struct {
	log       *slog.Logger
	libClient LibraryClient
}

type LibraryClient interface {
	Search(ctx context.Context, token jwt.Token, filter models.MediaFilter) ([]models.Media, error)
	NewMedia(ctx context.Context, token jwt.Token, media models.MediaConfig, source string) error
}

// TODO yandex client

func New(
	log *slog.Logger,
	libClient LibraryClient,
) *library {
	return &library{
		log:       log,
		libClient: libClient,
	}
}

func (l *library) Search(ctx context.Context, id int64, filter models.MediaFilter) ([]models.Media, error) {
	panic("not implemented")
}

func (l *library) NewMedia(ctx context.Context, id int64, media models.MediaConfig, source string) error {
	panic("not implemented")
}

func (l *library) LinkDownload(ctx context.Context, id int64, link string) (models.LinkDownloadResult, error) {
	panic("not implemented")
}

func (l *library) LinkUpload(ctx context.Context, id int64, res models.LinkDownloadResult) error {
	panic("not implemented")
}

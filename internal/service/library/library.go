package library

import (
	"context"
	"log/slog"

	"github.com/GintGld/fizteh-radio-bot/internal/models"
)

type library struct {
	log *slog.Logger
}

func New(
	log *slog.Logger,
) *library {
	return &library{
		log: log,
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

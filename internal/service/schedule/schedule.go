package schedule

import (
	"context"
	"log/slog"

	"github.com/GintGld/fizteh-radio-bot/internal/models"
)

type schedule struct {
	log *slog.Logger
}

func New(
	log *slog.Logger,
) *schedule {
	return &schedule{
		log: log,
	}
}

func (s *schedule) NewSegment(ctx context.Context, id int64, segm models.Segment) error {
	panic("not implemented")
}

func (s *schedule) Schedule(ctx context.Context, id int64) ([]models.Segment, error) {
	panic("not implemented")
}

func (s *schedule) Config(ctx context.Context, id int64) (models.AutoDJConfig, error) {
	panic("not implemented")
}

func (s *schedule) SetConfig(ctx context.Context, id int64, config models.AutoDJConfig) error {
	panic("not implemented")
}

func (s *schedule) Start(ctx context.Context) error {
	panic("not implemented")
}

func (s *schedule) Stop(ctx context.Context) error {
	panic("not implemented")
}

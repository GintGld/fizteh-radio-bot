package schedule

import (
	"context"
	"log/slog"

	"github.com/GintGld/fizteh-radio-bot/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type schedule struct {
	log       *slog.Logger
	schClient ScheduleClient
	djClient  AutoDJClient
}

type ScheduleClient interface {
	NewSegment(ctx context.Context, token jwt.Token, segm models.Segment) error
	GetSchedule(ctx context.Context, token jwt.Token) ([]models.Segment, error)
}

type AutoDJClient interface {
	GetConfig(ctx context.Context, token jwt.Token) (models.AutoDJConfig, error)
	StartAutoDJ(ctx context.Context, token jwt.Token) error
	StopAutoDJ(ctx context.Context, token jwt.Token) error
}

func New(
	log *slog.Logger,
	schClient ScheduleClient,
	djClient AutoDJClient,
) *schedule {
	return &schedule{
		log:       log,
		schClient: schClient,
		djClient:  djClient,
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

func (s *schedule) StartAutoDJ(ctx context.Context) error {
	panic("not implemented")
}

func (s *schedule) StopAutoDJ(ctx context.Context) error {
	panic("not implemented")
}

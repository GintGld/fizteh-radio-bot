package schedule

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/GintGld/fizteh-radio-bot/internal/lib/logger/sl"
	"github.com/GintGld/fizteh-radio-bot/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// It is supposed that no segment duration
	// will be less or equal to second.
	timeEps = time.Second
)

type schedule struct {
	log       *slog.Logger
	auth      Auth
	libClient LibraryClient
	schClient ScheduleClient
	djClient  AutoDJClient
}

type Auth interface {
	Token(ctx context.Context, id int64) (jwt.Token, error)
}

type LibraryClient interface {
	Media(ctx context.Context, token jwt.Token, id int64) (models.Media, error)
}

type ScheduleClient interface {
	NewSegment(ctx context.Context, token jwt.Token, segm models.Segment) error
	GetSchedule(ctx context.Context, token jwt.Token) ([]models.Segment, error)
}

type AutoDJClient interface {
	GetConfig(ctx context.Context, token jwt.Token) (models.AutoDJConfig, error)
	SetConfig(ctx context.Context, token jwt.Token, conf models.AutoDJConfig) error
	StartAutoDJ(ctx context.Context, token jwt.Token) error
	StopAutoDJ(ctx context.Context, token jwt.Token) error
	IsAutoDJPlaying(ctx context.Context, token jwt.Token) (bool, error)
}

func New(
	log *slog.Logger,
	auth Auth,
	libClient LibraryClient,
	schClient ScheduleClient,
	djClient AutoDJClient,
) *schedule {
	return &schedule{
		log:       log,
		auth:      auth,
		libClient: libClient,
		schClient: schClient,
		djClient:  djClient,
	}
}

func (s *schedule) NewSegment(ctx context.Context, id int64, segm models.Segment) error {
	const op = "schedule.NewSegment"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("userId", id),
	)

	token, err := s.auth.Token(ctx, id)
	if err != nil {
		log.Error(
			"failed to get token",
			sl.Err(err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.schClient.NewSegment(ctx, token, segm); err != nil {
		log.Error(
			"failed to create new segment",
			slog.String("name", segm.Media.Name),
			slog.String("start", segm.Start.String()),
			sl.Err(err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *schedule) AddToQueue(ctx context.Context, id int64, media models.MediaConfig) (models.Segment, error) {
	const op = "schedule.AddToQueue"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("userId", id),
		slog.Int64("mediaId", media.ID),
	)

	token, err := s.auth.Token(ctx, id)
	if err != nil {
		log.Error("failed to get token", sl.Err(err))
		return models.Segment{}, fmt.Errorf("%s: %w", op, err)
	}

	res, err := s.Schedule(ctx, id)
	if err != nil {
		log.Error("failed to get schedule", sl.Err(err))
		return models.Segment{}, fmt.Errorf("%s: %w", op, err)
	}

	now := time.Now()

	// Find first protected segment which is "after now".
	if j := slices.IndexFunc(res, func(s models.Segment) bool {
		return s.Start.Add(s.StopCut-s.BeginCut).After(now) && s.Protected
	}); j != -1 {
		// Cut uneccessary segments
		// (too early or not protected)
		res = slices.DeleteFunc(res[j:], func(s models.Segment) bool {
			return !s.Protected
		})
	} else {
		res = []models.Segment{}
	}

	supposedStart := now
	dur := media.Duration

	// If media does not fit into this time slot
	// move to the end of those segment.
	for _, segm := range res {
		if dur > segm.Start.Sub(supposedStart) {
			supposedStart = segm.Start.Add(segm.StopCut - segm.BeginCut + timeEps)
			continue
		}
		break
	}

	segm := models.Segment{
		Media:     media.ToMedia(),
		Start:     supposedStart,
		BeginCut:  0,
		StopCut:   media.Duration,
		Protected: true,
	}

	if err := s.schClient.NewSegment(ctx, token, segm); err != nil {
		log.Error("failed to add segment", sl.Err(err))
		return models.Segment{}, fmt.Errorf("%s: %w", op, err)
	}

	return segm, nil
}

func (s *schedule) Schedule(ctx context.Context, id int64) ([]models.Segment, error) {
	const op = "schedule.Schedule"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("userId", id),
	)

	token, err := s.auth.Token(ctx, id)
	if err != nil {
		log.Error(
			"failed to get token",
			sl.Err(err),
		)
		return []models.Segment{}, fmt.Errorf("%s: %w", op, err)
	}

	res, err := s.schClient.GetSchedule(ctx, token)
	if err != nil {
		log.Error(
			"failed to get schedule",
			sl.Err(err),
		)
		return []models.Segment{}, fmt.Errorf("%s: %w", op, err)
	}

	for i := range res {
		media, err := s.libClient.Media(ctx, token, res[i].Media.ID)
		if err != nil {
			log.Error(
				"failed to get media",
				slog.Int64("media id", res[i].Media.ID),
				sl.Err(err),
			)
			return []models.Segment{}, fmt.Errorf("%s: %w", op, err)
		}
		res[i].Media = media
	}

	return res, nil
}

func (s *schedule) Config(ctx context.Context, id int64) (models.AutoDJInfo, error) {
	const op = "schedule.Config"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("userId", id),
	)

	token, err := s.auth.Token(ctx, id)
	if err != nil {
		log.Error(
			"failed to get token",
			sl.Err(err),
		)
		return models.AutoDJInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	conf, err := s.djClient.GetConfig(ctx, token)
	if err != nil {
		log.Error(
			"failed to get schedule",
			sl.Err(err),
		)
		return models.AutoDJInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	info := conf.ToInfo()

	info.IsPlaying, err = s.djClient.IsAutoDJPlaying(ctx, token)
	if err != nil {
		log.Error(
			"failed to check if autodj playing",
		)
	}

	return info, nil
}

func (s *schedule) SetConfig(ctx context.Context, id int64, info models.AutoDJInfo) error {
	const op = "schedule.SetConfig"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("userId", id),
	)

	token, err := s.auth.Token(ctx, id)
	if err != nil {
		log.Error(
			"failed to get token",
			sl.Err(err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	conf := info.ToConfig()

	if err := s.djClient.SetConfig(ctx, token, conf); err != nil {
		log.Error(
			"failed to create new segment",
			sl.Err(err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *schedule) StartAutoDJ(ctx context.Context, id int64) error {
	const op = "schedule.StartAutoDJ"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("userId", id),
	)

	token, err := s.auth.Token(ctx, id)
	if err != nil {
		log.Error(
			"failed to get token",
			sl.Err(err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.djClient.StartAutoDJ(ctx, token); err != nil {
		log.Error(
			"failed to start autodj",
			sl.Err(err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *schedule) StopAutoDJ(ctx context.Context, id int64) error {
	const op = "schedule.StopAutoDJ"

	log := s.log.With(
		slog.String("op", op),
		slog.Int64("userId", id),
	)

	token, err := s.auth.Token(ctx, id)
	if err != nil {
		log.Error(
			"failed to get token",
			sl.Err(err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := s.djClient.StopAutoDJ(ctx, token); err != nil {
		log.Error(
			"failed to stop autodj",
			sl.Err(err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

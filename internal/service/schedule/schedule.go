package schedule

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/GintGld/fizteh-radio-bot/internal/lib/logger/sl"
	"github.com/GintGld/fizteh-radio-bot/internal/models"

	"github.com/golang-jwt/jwt/v5"
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

	var info models.AutoDJInfo

	for _, e := range conf.Tags {
		switch e.Type {
		case models.TagTypesAvail["genre"]:
			info.Genres = append(info.Genres, e.Name)
		case models.TagTypesAvail["playlist"]:
			info.Playlists = append(info.Playlists, e.Name)
		case models.TagTypesAvail["mood"]:
			info.Moods = append(info.Moods, e.Name)
		case models.TagTypesAvail["language"]:
			info.Languages = append(info.Languages, e.Name)
		}
	}

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

	tags := make(models.TagList, len(info.Genres)+len(info.Languages)+len(info.Moods)+len(info.Playlists))

	for _, e := range info.Genres {
		tags = append(tags, models.Tag{
			Name: e,
			Type: models.TagTypesAvail["genre"],
		})
	}
	for _, e := range info.Playlists {
		tags = append(tags, models.Tag{
			Name: e,
			Type: models.TagTypesAvail["playlist"],
		})
	}
	for _, e := range info.Moods {
		tags = append(tags, models.Tag{
			Name: e,
			Type: models.TagTypesAvail["mood"],
		})
	}
	for _, e := range info.Languages {
		tags = append(tags, models.Tag{
			Name: e,
			Type: models.TagTypesAvail["language"],
		})
	}

	conf := models.AutoDJConfig{
		Tags: tags,
	}

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

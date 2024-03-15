package library

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"

	"github.com/GintGld/fizteh-radio-bot/internal/lib/logger/sl"
	"github.com/GintGld/fizteh-radio-bot/internal/models"
	"github.com/GintGld/fizteh-radio-bot/internal/service"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: add support for spotify

var (
	yaSong     = regexp.MustCompile(`^https://music\.yandex\.(ru|com)/album/\d+/track/(?P<track>\d+)`)
	yaPlaylist = regexp.MustCompile(`^https://music\.yandex\.(ru|com)/users/(?P<user>[^\/]+)/playlists/(?P<track>\d+)`)
)

type library struct {
	log       *slog.Logger
	auth      Auth
	libClient LibraryClient
	yaClient  YaClient
}

type Auth interface {
	Token(ctx context.Context, id int64) (jwt.Token, error)
}

type LibraryClient interface {
	Search(ctx context.Context, token jwt.Token, filter models.MediaFilter) ([]models.Media, error)
	NewMedia(ctx context.Context, token jwt.Token, media models.Media) error
	NewTag(ctx context.Context, token jwt.Token, tag models.Tag) error
}

type YaClient interface {
	DownloadTrack(ctx context.Context, id int) (models.Media, error)
	DownloadPlaylist(ctx context.Context, user string, id int) (models.Playlist, error)
}

func New(
	log *slog.Logger,
	auth Auth,
	libClient LibraryClient,
	yaClient YaClient,
) *library {
	l := &library{
		log:       log,
		auth:      auth,
		libClient: libClient,
		yaClient:  yaClient,
	}

	return l
}

func (l *library) Search(ctx context.Context, id int64, filter models.MediaFilter) ([]models.Media, error) {
	const op = "library.Search"

	log := l.log.With(
		slog.String("op", op),
	)

	token, err := l.auth.Token(ctx, id)
	if err != nil {
		log.Error(
			"failed to get token",
			slog.Int64("id", id),
			sl.Err(err),
		)
		return []models.Media{}, fmt.Errorf("%s: %w", op, err)
	}

	res, err := l.libClient.Search(ctx, token, filter)
	if err != nil {
		log.Error(
			"failed to get library",
			slog.Int64("id", id),
			sl.Err(err),
		)
		return []models.Media{}, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (l *library) NewMedia(ctx context.Context, id int64, mediaConf models.MediaConfig, source string) error {
	const op = "NewMedia"

	log := l.log.With(
		slog.String("op", op),
	)

	token, err := l.auth.Token(ctx, id)
	if err != nil {
		log.Error(
			"failed to get token",
			slog.Int64("id", id),
			sl.Err(err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	tags := make(models.TagList,
		len(mediaConf.Playlists)+
			len(mediaConf.Podcasts)+
			len(mediaConf.Genres)+
			len(mediaConf.Languages)+
			len(mediaConf.Moods),
	)
	for i := range mediaConf.Playlists {
		tags[i] = models.Tag{}
	}

	media := models.Media{
		Name:       mediaConf.Name,
		Author:     mediaConf.Author,
		Duration:   mediaConf.Duration,
		SourcePath: source,
	}

	if err := l.libClient.NewMedia(ctx, token, media); err != nil {
		log.Error(
			"failed to upload media",
			slog.Int64("id", id),
			sl.Err(err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (l *library) LinkDownload(ctx context.Context, id int64, link string) (models.LinkDownloadResult, error) {
	const op = "library.LinkDownload"

	log := l.log.With(
		slog.String("op", op),
	)

	var res models.LinkDownloadResult

	if yaSong.MatchString(link) {
		res.Type = models.ResSong

		subStr := string(yaSong.ExpandString([]byte{}, "$track", link, yaSong.FindSubmatchIndex([]byte(link))))

		trackId, err := strconv.Atoi(subStr)
		if err != nil {
			log.Error("failed to convert id to int",
				slog.String("link", link),
				slog.String("expected id", subStr),
				sl.Err(err),
			)
			return models.LinkDownloadResult{}, fmt.Errorf("%s: %w", op, err)
		}

		media, err := l.yaClient.DownloadTrack(ctx, trackId)
		if err != nil {
			// TODO handler errors
			log.Error(
				"failed to download track from yandex",
				slog.Int("trackId", trackId),
				sl.Err(err),
			)
			return models.LinkDownloadResult{}, fmt.Errorf("%s: %w", op, err)
		}

		res.Media = media
	} else if yaPlaylist.MatchString(link) {
		res.Type = models.ResPlaylist

		userName := string(yaPlaylist.ExpandString([]byte{}, "$user", link, yaPlaylist.FindSubmatchIndex([]byte(link))))

		subStr := string(yaSong.ExpandString([]byte{}, "$track", link, yaSong.FindSubmatchIndex([]byte(link))))
		id, err := strconv.Atoi(subStr)
		if err != nil {
			log.Error("failed to convert id to int",
				slog.String("link", link),
				slog.String("expected id", subStr),
				sl.Err(err),
			)
			return models.LinkDownloadResult{}, fmt.Errorf("%s: %w", op, err)
		}

		playlist, err := l.yaClient.DownloadPlaylist(ctx, userName, id)
		if err != nil {
			// TODO handler errors
			log.Error(
				"failed to download track from yandex",
				slog.String("user", userName),
				slog.Int("playlist", id),
				sl.Err(err),
			)
			return models.LinkDownloadResult{}, fmt.Errorf("%s: %w", op, err)
		}

		res.Playlist = playlist
	} else {
		log.Warn(
			"unknown link",
			slog.String("link", link),
		)
		return models.LinkDownloadResult{}, service.ErrInvalidLink
	}

	return res, nil
}

func (l *library) LinkUpload(ctx context.Context, id int64, res models.LinkDownloadResult) error {
	const op = "library.Search"

	log := l.log.With(
		slog.String("op", op),
	)

	token, err := l.auth.Token(ctx, id)
	if err != nil {
		log.Error(
			"failed to get token",
			slog.Int64("id", id),
			sl.Err(err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	switch res.Type {
	case models.ResSong:
		l.libClient.NewMedia(ctx, token, res.Media)
	case models.ResPlaylist:
		if err := l.libClient.NewTag(ctx, token, models.Tag{
			Name: res.Playlist.Name,
			Type: models.TagTypesAvail[3], // tag type "playlist"
		}); err != nil {
			// TODO: handle "tag already exist" tag.
			log.Error(
				"failed to create tag",
				slog.String("name", res.Playlist.Name),
				sl.Err(err),
			)
			return fmt.Errorf("%s: %w", op, err)
		}

		for _, m := range res.Playlist.Values {
			if err := l.libClient.NewMedia(ctx, token, m); err != nil {
				log.Error(
					"failed to upload media",
					slog.String("Name", m.Name),
					slog.String("source", m.SourcePath),
					sl.Err(err),
				)
			}
		}
	}

	return nil
}

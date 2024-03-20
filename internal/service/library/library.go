package library

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"slices"
	"strconv"
	"time"

	"github.com/GintGld/fizteh-radio-bot/internal/client"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/logger/sl"
	"github.com/GintGld/fizteh-radio-bot/internal/models"
	yamodels "github.com/GintGld/fizteh-radio-bot/internal/models/yandex"
	"github.com/GintGld/fizteh-radio-bot/internal/service"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: add support for spotify

var (
	yaSong     = regexp.MustCompile(`^https://music\.yandex\.(ru|com)/album/(?P<album>\d+)/track/(?P<track>\d+)`)
	yaPlaylist = regexp.MustCompile(`^https://music\.yandex\.(ru|com)/users/(?P<user>[^\/]+)/playlists/(?P<kind>\d+)`)
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
	Album(ctx context.Context, id string) (yamodels.Album, error)
	Playlist(ctx context.Context, user string, id string) (yamodels.Playlist, error)
	DownloadInfo(ctx context.Context, id int) ([]yamodels.DownloadInfo, error)
	DownloadTrack(ctx context.Context, url string) (string, error)
	DirectLink(ctx context.Context, url string) (string, error)
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
		slog.Int64("userId", id),
	)

	token, err := l.auth.Token(ctx, id)
	if err != nil {
		log.Error(
			"failed to get token",
			sl.Err(err),
		)
		return []models.Media{}, fmt.Errorf("%s: %w", op, err)
	}

	res, err := l.libClient.Search(ctx, token, filter)
	if err != nil {
		log.Error(
			"failed to get library",
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
		slog.Int64("userId", id),
	)

	token, err := l.auth.Token(ctx, id)
	if err != nil {
		log.Error(
			"failed to get token",
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
	for _, t := range mediaConf.Playlists {
		tags = append(tags, models.Tag{
			Name: t,
			Type: models.TagTypesAvail["playlist"],
		})
	}
	for _, t := range mediaConf.Podcasts {
		tags = append(tags, models.Tag{
			Name: t,
			Type: models.TagTypesAvail["podcast"],
		})
	}
	for _, t := range mediaConf.Genres {
		tags = append(tags, models.Tag{
			Name: t,
			Type: models.TagTypesAvail["genre"],
		})
	}
	for _, t := range mediaConf.Languages {
		tags = append(tags, models.Tag{
			Name: t,
			Type: models.TagTypesAvail["language"],
		})
	}
	for _, t := range mediaConf.Moods {
		tags = append(tags, models.Tag{
			Name: t,
			Type: models.TagTypesAvail["mood"],
		})
	}

	media := models.Media{
		Name:       mediaConf.Name,
		Author:     mediaConf.Author,
		Duration:   mediaConf.Duration,
		Tags:       tags,
		SourcePath: source,
	}

	if err := l.libClient.NewMedia(ctx, token, media); err != nil {
		log.Error(
			"failed to upload media",
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
		slog.Int64("userId", id),
	)

	var res models.LinkDownloadResult

	if yaSong.MatchString(link) {
		res.Type = models.ResSong

		media, err := l.linkTrack(ctx, link)
		if err != nil {
			log.Error(
				"failed to handle track",
				sl.Err(err),
			)
			return models.LinkDownloadResult{}, fmt.Errorf("%s: %w", op, err)
		}

		res.Media = media
	} else if yaPlaylist.MatchString(link) {
		res.Type = models.ResPlaylist

		playlist, err := l.linkPlaylist(ctx, link)
		if err != nil {
			log.Error(
				"failed to handle playlist",
				sl.Err(err),
			)
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

func (l *library) linkTrack(ctx context.Context, url string) (models.Media, error) {
	const op = "library.linkTrack"

	log := l.log.With(
		slog.String("op", op),
	)

	trackId := string(yaSong.ExpandString([]byte{}, "$track", url, yaSong.FindSubmatchIndex([]byte(url))))
	albumId := string(yaSong.ExpandString([]byte{}, "$album", url, yaSong.FindSubmatchIndex([]byte(url))))

	log.Debug(
		"parsed values",
		slog.String("trackId", trackId),
		slog.String("albumId", albumId),
	)

	album, err := l.yaClient.Album(ctx, albumId)
	if err != nil {
		// TODO handle errors
		log.Error(
			"failed to get album info",
			slog.String("albumId", albumId),
			sl.Err(err),
		)
		return models.Media{}, fmt.Errorf("%s: %w", op, err)
	}
	if album.Err != nil {
		log.Error(
			"failed to get album info",
			slog.String("albumId", albumId),
			slog.String("error", *album.Err),
		)
		return models.Media{}, fmt.Errorf("%s: %w", op, err)
	}

	index := slices.IndexFunc(album.Tracks, func(t yamodels.Track) bool {
		return strconv.Itoa(t.Id) == trackId
	})
	track := album.Tracks[index]

	filePath, err := l.downloadLinkTrack(ctx, track.Id)
	if err != nil {
		// TODO handler errors
		log.Error(
			"failed to download track from yandex",
			slog.Int("trackId", track.Id),
			sl.Err(err),
		)
		return models.Media{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.Media{
		Name:       track.Title,
		Author:     track.Artists[0].Name,
		Duration:   track.Duration,
		SourcePath: filePath,
	}, nil
}

func (l *library) linkPlaylist(ctx context.Context, url string) (models.Playlist, error) {
	const op = "library.linkPlaylist"

	log := l.log.With(
		slog.String("op", op),
	)

	userName := string(yaPlaylist.ExpandString([]byte{}, "$user", url, yaPlaylist.FindSubmatchIndex([]byte(url))))
	kind := string(yaPlaylist.ExpandString([]byte{}, "$kind", url, yaSong.FindSubmatchIndex([]byte(url))))

	log.Debug(
		"parsed values",
		slog.String("user", userName),
		slog.String("kind", kind),
	)

	playlist, err := l.yaClient.Playlist(ctx, userName, kind)
	if err != nil {
		// TODO handler errors
		log.Error(
			"failed to download track from yandex",
			slog.String("user", userName),
			slog.String("kind", kind),
			sl.Err(err),
		)
		return models.Playlist{}, fmt.Errorf("%s: %w", op, err)
	}

	values := make([]models.Media, 0, len(playlist.Tracks))
	for _, track := range playlist.Tracks {
		filePath, err := l.downloadLinkTrack(ctx, track.Id)
		if err != nil {
			log.Error(
				"failed to download track",
				slog.Int("trackId", track.Id),
				sl.Err(err),
			)
			return models.Playlist{}, fmt.Errorf("%s: %w", op, err)
		}

		values = append(values, models.Media{
			Name:       track.Title,
			Author:     track.Artists[0].Name,
			Duration:   track.Duration,
			SourcePath: filePath,
		})
	}

	return models.Playlist{
		Name:   playlist.Title,
		Values: values,
	}, nil
}

// downloadLinkTrack downloads
// track file and returns path to the file.
func (l *library) downloadLinkTrack(ctx context.Context, id int) (string, error) {
	const op = "library.downloadLinkTrack"

	log := l.log.With(
		slog.String("op", op),
		slog.Int("tracakId", id),
	)

	downloadOptions, err := l.yaClient.DownloadInfo(ctx, id)
	if err != nil {
		log.Error(
			"failed to get download options",
			sl.Err(err),
		)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	downloadOptions = slices.DeleteFunc(downloadOptions, func(di yamodels.DownloadInfo) bool {
		return di.Codec != yamodels.CodecMP3
	})

	if len(downloadOptions) == 0 {
		log.Warn(
			"download info with mp3 codec not found",
		)
		return "", client.ErrTrackNotFound
	}

	preferred := slices.MaxFunc(downloadOptions, func(a, b yamodels.DownloadInfo) int {
		return int(b.Bitrate - a.Bitrate)
	})

	directURL, err := l.yaClient.DirectLink(ctx, preferred.URL)
	if err != nil {
		log.Error(
			"failed to get direct download link",
			sl.Err(err),
		)
		return "", client.ErrTrackNotFound
	}

	filePath, err := l.yaClient.DownloadTrack(ctx, directURL)
	if err != nil {
		log.Error(
			"failed to download track",
			sl.Err(err),
		)
		return "", client.ErrTrackNotFound
	}

	time.AfterFunc(time.Hour, func() {
		if err := os.Remove(filePath); err != nil {
			l.log.Error(
				"failed to delete file",
				slog.String("path", filePath),
				sl.Err(err),
			)
		}
	})

	return filePath, nil
}

func (l *library) LinkUpload(ctx context.Context, id int64, res models.LinkDownloadResult) error {
	const op = "library.Search"

	log := l.log.With(
		slog.String("op", op),
		slog.Int64("userId", id),
	)

	token, err := l.auth.Token(ctx, id)
	if err != nil {
		log.Error(
			"failed to get token",
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
			Type: models.TagTypesAvail["playlists"],
		}); err != nil {
			// TODO: handle "tag already exist".
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

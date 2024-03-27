package library

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"slices"
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
	Media(ctx context.Context, token jwt.Token, id int64) (models.Media, error)
	NewMedia(ctx context.Context, token jwt.Token, media models.Media) (int64, error)
	UpdateMedia(ctx context.Context, token jwt.Token, media models.Media) error
	DeleteMedia(ctx context.Context, token jwt.Token, mediaId int64) error
	AllTags(ctx context.Context, token jwt.Token) (models.TagList, error)
	NewTag(ctx context.Context, token jwt.Token, tag models.Tag) (int64, error)
}

type YaClient interface {
	Album(ctx context.Context, id string) (yamodels.Album, error)
	Playlist(ctx context.Context, user string, id string) (yamodels.Playlist, error)
	DownloadInfo(ctx context.Context, id string) ([]yamodels.DownloadInfo, error)
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

func (l *library) Search(ctx context.Context, id int64, filter models.MediaFilter) ([]models.MediaConfig, error) {
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
		return []models.MediaConfig{}, fmt.Errorf("%s: %w", op, err)
	}

	res, err := l.libClient.Search(ctx, token, filter)
	if err != nil {
		log.Error(
			"failed to get library",
			sl.Err(err),
		)
		return []models.MediaConfig{}, fmt.Errorf("%s: %w", op, err)
	}

	configs := make([]models.MediaConfig, 0, len(res))
	for _, m := range res {
		configs = append(configs, m.ToConfig())
	}

	return configs, nil
}

func (l *library) NewMedia(ctx context.Context, id int64, mediaConf models.MediaConfig) (int64, error) {
	const op = "library.NewMedia"

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
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	media := mediaConf.ToMedia()

	if err := l.recoverTagIds(ctx, token, &media); err != nil {
		if errors.Is(err, service.ErrTagNotFound) {
			log.Warn("media has not existing tag", slog.String("name", media.Name), slog.Any("tags", media.Tags))
			return 0, service.ErrTagNotFound
		}
	}

	// search media with similar names and authors.
	searchRes, err := l.Search(ctx, id, models.MediaFilter{
		Name:       media.Name,
		Author:     media.Author,
		MaxRespLen: 20,
	})
	if err != nil {
		log.Error("failed to make search", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Check if media already exists.
	if index := slices.IndexFunc(searchRes, func(conf models.MediaConfig) bool {
		return conf.Name == media.Name && conf.Author == media.Author
	}); index != -1 {
		log.Warn("media already exists", slog.String("name", media.Name), slog.String("author", media.Author))
		return searchRes[index].ID, service.ErrMediaExists
	}

	mediaId, err := l.libClient.NewMedia(ctx, token, media)
	if err != nil {
		log.Error(
			"failed to upload media",
			sl.Err(err),
		)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return mediaId, nil
}

func (l *library) UpdateMedia(ctx context.Context, id int64, mediaConf models.MediaConfig) error {
	const op = "library.UpdateMedia"

	log := l.log.With(
		slog.String("op", op),
		slog.Int64("userId", id),
		slog.Int64("mediaId", mediaConf.ID),
	)

	token, err := l.auth.Token(ctx, id)
	if err != nil {
		log.Error(
			"failed to get token",
			sl.Err(err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	media := mediaConf.ToMedia()

	tagAvail, err := l.libClient.AllTags(ctx, token)
	if err != nil {
		log.Error("failed to get available tags", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	for i, tag := range media.Tags {
		if j := slices.IndexFunc(tagAvail, func(t models.Tag) bool {
			return t.Name == tag.Name
		}); j != -1 {
			log.Debug("found", slog.Int("j", j))
			media.Tags[i] = tagAvail[j]
		} else {
			log.Debug("create")
			id, err := l.libClient.NewTag(ctx, token, tag)
			if err != nil {
				log.Error(
					"failed to create tag",
					slog.String("name", tag.Name),
					sl.Err(err),
				)
				return fmt.Errorf("%s: %w", op, err)
			}
			media.Tags[i].ID = id
		}
	}

	if err := l.libClient.UpdateMedia(ctx, token, media); err != nil {
		log.Error(
			"failed to update media",
			sl.Err(err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (l *library) DeleteMedia(ctx context.Context, id int64, mediaConf models.MediaConfig) error {
	const op = "library.DeleteMedia"

	log := l.log.With(
		slog.String("op", op),
		slog.Int64("userId", id),
		slog.Int64("mediaId", mediaConf.ID),
	)

	token, err := l.auth.Token(ctx, id)
	if err != nil {
		log.Error(
			"failed to get token",
			sl.Err(err),
		)
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := l.libClient.DeleteMedia(ctx, token, mediaConf.ID); err != nil {
		log.Error(
			"failed to delete media",
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

		mediaConf, err := l.linkTrack(ctx, id, link)
		if err != nil {
			log.Error(
				"failed to handle track",
				sl.Err(err),
			)
			return models.LinkDownloadResult{}, fmt.Errorf("%s: %w", op, err)
		}

		res.MediaConf = mediaConf
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

func (l *library) linkTrack(ctx context.Context, id int64, url string) (models.MediaConfig, error) {
	const op = "library.linkTrack"

	log := l.log.With(
		slog.String("op", op),
		slog.Int64("userId", id),
	)

	trackId, albumId := exctractTrackInfo(url)

	album, err := l.yaClient.Album(ctx, albumId)
	if err != nil {
		// TODO handle errors
		log.Error(
			"failed to get album info",
			slog.String("albumId", albumId),
			sl.Err(err),
		)
		return models.MediaConfig{}, fmt.Errorf("%s: %w", op, err)
	}
	if album.Err != nil {
		log.Error(
			"failed to get album info",
			slog.String("albumId", albumId),
			slog.String("error", *album.Err),
		)
		return models.MediaConfig{}, fmt.Errorf("%s: %w", op, err)
	}

	index := slices.IndexFunc(album.Tracks, func(t yamodels.Track) bool {
		return t.Id == trackId
	})
	track := album.Tracks[index]

	filePath, err := l.downloadLinkTrack(ctx, track.Id)
	if err != nil {
		// TODO handler errors
		log.Error(
			"failed to download track from yandex",
			slog.String("trackId", track.Id),
			sl.Err(err),
		)
		return models.MediaConfig{}, fmt.Errorf("%s: %w", op, err)
	}

	var tag models.MediaFormat

	switch track.Format {
	case yamodels.YaMusicFormat:
		tag = models.Song
	case yamodels.YaPodcastFormat:
		tag = models.Podcast
	}

	author := ""
	if len(track.Artists) > 0 {
		author = track.Artists[0].Name
	}

	return models.MediaConfig{
		Name:       track.Title,
		Author:     author,
		Duration:   track.Duration,
		SourcePath: filePath,
		Format:     tag,
	}, nil
}

func (l *library) linkPlaylist(ctx context.Context, url string) (models.Playlist, error) {
	const op = "library.linkPlaylist"

	log := l.log.With(
		slog.String("op", op),
	)

	kind, userName := exctractPlaylistInfo(url)

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

	values := make([]models.MediaConfig, 0, len(playlist.Tracks))
	for _, track := range playlist.Tracks {
		filePath, err := l.downloadLinkTrack(ctx, track.Id)
		if err != nil {
			log.Error(
				"failed to download track",
				slog.String("trackId", track.Id),
				sl.Err(err),
			)
			return models.Playlist{}, fmt.Errorf("%s: %w", op, err)
		}

		values = append(values, models.MediaConfig{
			Name:       track.Title,
			Author:     track.Artists[0].Name,
			Duration:   track.Duration,
			SourcePath: filePath,
			Format:     models.Song,
			Playlists:  []string{playlist.Title},
		})
	}

	return models.Playlist{
		Name:   playlist.Title,
		Values: values,
	}, nil
}

// downloadLinkTrack downloads
// track file and returns path to the file.
func (l *library) downloadLinkTrack(ctx context.Context, id string) (string, error) {
	const op = "library.downloadLinkTrack"

	log := l.log.With(
		slog.String("op", op),
		slog.String("tracakId", id),
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
	const op = "library.LinkUpload"

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

	media := res.MediaConf.ToMedia()

	switch res.Type {
	case models.ResSong:
		if _, err := l.libClient.NewMedia(ctx, token, media); err != nil {
			if errors.Is(err, service.ErrMediaExists) {
				log.Info(
					"media already exists",
					slog.String("name", media.Name),
					slog.String("author", media.Author),
				)
				return service.ErrMediaExists
			}
			log.Error(
				"failed to upload media",
				slog.String("name", media.Name),
				sl.Err(err),
			)
			return fmt.Errorf("%s: %w", op, err)
		}
	case models.ResPlaylist:
		if _, err := l.libClient.NewTag(ctx, token, models.Tag{
			Name: res.Playlist.Name,
			Type: models.TagTypesAvail["playlist"],
		}); err != nil {
			// TODO: handle "tag already exist".
			log.Error(
				"failed to create tag",
				slog.String("name", res.Playlist.Name),
				sl.Err(err),
			)
			return fmt.Errorf("%s: %w", op, err)
		}

	upload_loop:
		for _, m := range res.Playlist.Values {
			if id, err := l.NewMedia(ctx, id, m); err != nil {
				// If media already exists, add new tags to it
				if errors.Is(err, service.ErrMediaExists) {
					log.Info(
						"media already exists, add new tags to it",
						slog.String("name", m.Name),
						slog.String("author", m.Author),
					)
					m := m.ToMedia()
					oldMedia, err := l.libClient.Media(ctx, token, id)
					if err != nil {
						log.Error("failed to get media", slog.Int64("id", id), sl.Err(err))
						continue upload_loop
					}
					// add new tags
					newTags := slices.Clone(oldMedia.Tags)
					for _, t := range m.Tags {
						if !slices.ContainsFunc(oldMedia.Tags, func(tInner models.Tag) bool {
							return tInner.Name == t.Name && tInner.Type == t.Type
						}) {
							newTags = append(newTags, t)
						}
					}
					oldMedia.Tags = newTags
					if err := l.libClient.UpdateMedia(ctx, token, oldMedia); err != nil {
						log.Error("failed to update media", slog.Int64("id", oldMedia.ID), sl.Err(err))
					}
					continue upload_loop
				}
				log.Error(
					"failed to upload media",
					slog.String("Name", m.Name),
					slog.String("source", m.SourcePath),
					sl.Err(err),
				)
				return fmt.Errorf("%s: %w", op, err)
			}
		}
	}

	return nil
}

func (l *library) recoverTagIds(ctx context.Context, token jwt.Token, m *models.Media) error {
	const op = "library.recoverTagIds"

	log := l.log.With(
		slog.String("op", op),
	)

	tags, err := l.libClient.AllTags(ctx, token)
	if err != nil {
		log.Error("failedt to get all tags")
		return fmt.Errorf("%s: %w", op, err)
	}

	for i, t := range m.Tags {
		if j := slices.IndexFunc(tags, func(tInner models.Tag) bool {
			return t.Name == tInner.Name && t.Type == tInner.Type
		}); j != -1 {
			m.Tags[i].ID = tags[j].ID
		} else {
			return service.ErrTagNotFound
		}
	}
	return nil
}

func exctractTrackInfo(url string) (trackId, albumId string) {
	trackId = string(yaSong.ExpandString([]byte{}, "$track", url, yaSong.FindSubmatchIndex([]byte(url))))
	albumId = string(yaSong.ExpandString([]byte{}, "$album", url, yaSong.FindSubmatchIndex([]byte(url))))

	return
}

func exctractPlaylistInfo(url string) (kind, user string) {
	kind = string(yaPlaylist.ExpandString([]byte{}, "$kind", url, yaPlaylist.FindSubmatchIndex([]byte(url))))
	user = string(yaPlaylist.ExpandString([]byte{}, "$user", url, yaPlaylist.FindSubmatchIndex([]byte(url))))

	return
}

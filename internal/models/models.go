package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Login string    `json:"login"`
	Pass  string    `json:"pass"`
	Token jwt.Token `json:"-"`
}

type Media struct {
	ID         int64         `json:"id"`
	Name       string        `json:"name"`
	Author     string        `json:"author"`
	Duration   time.Duration `json:"duration"`
	Tags       TagList       `json:"tags"`
	SourcePath string        `json:"-"`
}

type Playlist struct {
	Name   string
	Values []MediaConfig
}

type MediaConfig struct {
	ID         int64
	Name       string
	Author     string
	Duration   time.Duration
	Format     MediaFormat
	Playlists  []string
	Podcasts   []string
	Genres     []string
	Languages  []string
	Moods      []string
	SourcePath string
}

type MediaFormat int

const (
	Song MediaFormat = iota
	Podcast
	Jingle
)

func (m MediaFormat) String() string {
	switch m {
	case Song:
		return "песня"
	case Podcast:
		return "подкаст"
	case Jingle:
		return "джингл"
	default:
		return ""
	}
}

type LinkDownloadResult struct {
	Type      ResultType
	MediaConf MediaConfig
	Playlist  Playlist
}

type ResultType int

const (
	ResSong ResultType = iota
	ResPlaylist
)

type MediaFilter struct {
	Name       string
	Author     string
	Tags       []string
	MaxRespLen int
}

type AutoDJInfo struct {
	IsPlaying bool
	Genres    []string
	Playlists []string
	Languages []string
	Moods     []string
}

// TODO add stub
type AutoDJConfig struct {
	Tags TagList `json:"tags"`
}

type TagTypes []TagType
type TagList []Tag

type Tag struct {
	ID   int64   `json:"id"`
	Name string  `json:"name"`
	Type TagType `json:"type"`
}

type TagType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

var (
	TagTypesAvail = map[string]TagType{
		"format":   {ID: 1, Name: "format"},
		"genre":    {ID: 2, Name: "genre"},
		"playlist": {ID: 3, Name: "playlist"},
		"mood":     {ID: 4, Name: "mood"},
		"language": {ID: 5, Name: "language"},
		"podcast":  {ID: 6, Name: "podcast"},
	}
)

type Segment struct {
	ID        int64
	Media     Media
	Start     time.Time
	BeginCut  time.Duration
	StopCut   time.Duration
	Protected bool
}

func (s *Segment) UnmarshalJSON(data []byte) error {
	var tmp segmentResponse

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	*s = Segment{
		ID: tmp.ID,
		Media: Media{
			ID: tmp.MediaID,
		},
		Start:     tmp.Start,
		BeginCut:  tmp.BeginCut,
		StopCut:   tmp.StopCut,
		Protected: tmp.Protected,
	}

	return nil
}

type segmentResponse struct {
	ID        int64         `json:"id"`
	MediaID   int64         `json:"mediaID"`
	Start     time.Time     `json:"start"`
	BeginCut  time.Duration `json:"beginCut"`
	StopCut   time.Duration `json:"stopCut"`
	Protected bool          `json:"protected"`
}

func (conf MediaConfig) ToMedia() Media {
	tags := make(TagList, 0,
		1+len(conf.Playlists)+
			len(conf.Podcasts)+
			len(conf.Genres)+
			len(conf.Languages)+
			len(conf.Moods),
	)
	switch conf.Format {
	case Song:
		tags = append(tags, Tag{
			Name: "song",
			Type: TagTypesAvail["format"],
		})
	case Podcast:
		tags = append(tags, Tag{
			Name: "podcast",
			Type: TagTypesAvail["format"],
		})
	case Jingle:
		tags = append(tags, Tag{
			Name: "jingle",
			Type: TagTypesAvail["format"],
		})
	}
	for _, t := range conf.Playlists {
		tags = append(tags, Tag{
			Name: t,
			Type: TagTypesAvail["playlist"],
		})
	}
	for _, t := range conf.Podcasts {
		tags = append(tags, Tag{
			Name: t,
			Type: TagTypesAvail["podcast"],
		})
	}
	for _, t := range conf.Genres {
		tags = append(tags, Tag{
			Name: t,
			Type: TagTypesAvail["genre"],
		})
	}
	for _, t := range conf.Languages {
		tags = append(tags, Tag{
			Name: t,
			Type: TagTypesAvail["language"],
		})
	}
	for _, t := range conf.Moods {
		tags = append(tags, Tag{
			Name: t,
			Type: TagTypesAvail["mood"],
		})
	}
	return Media{
		ID:         conf.ID,
		Name:       conf.Name,
		Author:     conf.Author,
		Duration:   conf.Duration,
		Tags:       tags,
		SourcePath: conf.SourcePath,
	}
}

func (m Media) ToConfig() MediaConfig {
	Playlists := make([]string, 0)
	Podcasts := make([]string, 0)
	Genres := make([]string, 0)
	Languages := make([]string, 0)
	Moods := make([]string, 0)
	var format MediaFormat

	for _, t := range m.Tags {
		switch t.Type.Name {
		case "format":
			switch t.Name {
			case "song":
				format = Song
			case "podcast":
				format = Podcast
			case "jingle":
				format = Jingle
			}
		case "playlist":
			Playlists = append(Playlists, t.Name)
		case "podcast":
			Podcasts = append(Podcasts, t.Name)
		case "genre":
			Genres = append(Genres, t.Name)
		case "language":
			Languages = append(Languages, t.Name)
		case "mood":
			Moods = append(Moods, t.Name)
		}
	}

	return MediaConfig{
		ID:        m.ID,
		Name:      m.Name,
		Author:    m.Author,
		Duration:  m.Duration,
		Format:    format,
		Playlists: Playlists,
		Podcasts:  Podcasts,
		Genres:    Genres,
		Languages: Languages,
		Moods:     Moods,
	}
}

func (conf MediaConfig) String() string {
	var b strings.Builder

	b.WriteString("<b>Композиция</b>\n")
	b.WriteString(fmt.Sprintf("<b>Название:</b> %s\n", conf.Name))
	b.WriteString(fmt.Sprintf("<b>Автор:</b> %s\n", conf.Author))
	b.WriteString(fmt.Sprintf("<b>Формат:</b> %s\n", conf.Format))
	b.WriteString(fmt.Sprintf("<b>Длительность:</b> %s\n", conf.Duration.Round(time.Second).String()))

	if len(conf.Podcasts) > 0 {
		b.WriteString(fmt.Sprintf("<b>Подкасты:</b> %s\n", strings.Join(conf.Podcasts, ", ")))
	}
	if len(conf.Playlists) > 0 {
		b.WriteString(fmt.Sprintf("<b>Плейлисты:</b> %s\n", strings.Join(conf.Playlists, ", ")))
	}
	if len(conf.Genres) > 0 {
		b.WriteString(fmt.Sprintf("<b>Жанры:</b> %s\n", strings.Join(conf.Genres, ", ")))
	}
	if len(conf.Languages) > 0 {
		b.WriteString(fmt.Sprintf("<b>Языки:</b> %s\n", strings.Join(conf.Languages, ", ")))
	}
	if len(conf.Moods) > 0 {
		b.WriteString(fmt.Sprintf("<b>Настроение:</b> %s\n", strings.Join(conf.Moods, ", ")))
	}

	return b.String()
}

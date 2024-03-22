package models

import (
	"encoding/json"
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
	Values []Media
}

type MediaConfig struct {
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
)

type LinkDownloadResult struct {
	Type     ResultType
	Media    Media
	Playlist Playlist
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

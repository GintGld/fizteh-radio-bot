package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// type Editor struct {
// 	Login string `json:"login"`
// 	Pass  string `json:"pass"`
// }

type User struct {
	Login string
	Pass  string
	Token jwt.Token
}

type Media struct {
	ID       int64         `json:"id"`
	Name     string        `json:"name"`
	Author   string        `json:"author"`
	Duration time.Duration `json:"duration"`
	Tags     TagList       `json:"tags"`
}

type MediaConfig struct {
	Name      string
	Author    string
	Duration  time.Duration
	Format    MediaFormat
	Playlists []string
	Podcasts  []string
	Genres    []string
	Languages []string
	Moods     []string
}

type MediaFormat int

const (
	Song MediaFormat = iota
	Podcast
)

type LinkDownloadResult struct {
	Type     ResultType
	Media    Media
	Playlist struct {
		Name   string
		Values []Media
	}
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

type AutoDJConfig struct {
	IsPlaying bool
	Genres    []string
	Playlists []string
	Languages []string
	Moods     []string
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
	TagTypesAvail = TagTypes{
		TagType{ID: 1, Name: "format"},
		TagType{ID: 2, Name: "genre"},
		TagType{ID: 3, Name: "playlist"},
		TagType{ID: 4, Name: "mood"},
		TagType{ID: 5, Name: "language"},
	}
)

type Segment struct {
	ID        int64         `json:"id"`
	Media     Media         `json:"mediaID"`
	Start     time.Time     `json:"start"`
	BeginCut  time.Duration `json:"beginCut"`
	StopCut   time.Duration `json:"stopCut"`
	Protected bool          `json:"protected"`
}

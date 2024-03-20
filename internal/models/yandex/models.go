package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"time"
)

const (
	SIGN_SALT = "XGRlBW9FXlekgbPrRHuSiA"
)

var (
	ErrTrackIsNotMusic = errors.New("track is not music")
)

type DownloadInfo struct {
	Codec   CodecType `json:"codec"`
	Gain    bool      `json:"gain"`
	Preview string    `json:"preview"`
	URL     string    `json:"downloadInfoUrl"`
	Direct  bool      `json:"direct"`
	Bitrate float64   `json:"bitrateInKbps"`
}

type CodecType string

const (
	CodecMP3 = "mp3"
	CodeAAC  = "aac"
)

type DownloadInfoXMLResponse struct {
	XMLName xml.Name `xml:"download-info"`
	Host    string   `xml:"host"`
	Path    string   `xml:"path"`
	Ts      string   `xml:"ts"`
	S       string   `xml:"s"`
}

func (d DownloadInfoXMLResponse) BuildLink() string {
	h := md5.New()
	h.Write([]byte(SIGN_SALT + d.Path[1:] + d.S))
	sign := hex.EncodeToString(h.Sum(nil))

	// 'https://{host}/get-mp3/{sign}/{ts}{path}'
	return fmt.Sprintf("https://%s/get-mp3/%s/%s%s", d.Host, sign, d.Ts, d.Path)
}

type Playlist struct {
	Id       int
	Kind     int
	Title    string
	Duration time.Duration
	Tracks   []Track
}

type Album struct {
	Id       int
	Err      *string
	Title    string
	MetaType MetaType
	Genre    string
	Artists  []Artist
	Tracks   []Track
}

type MetaType string

const (
	Single  MetaType = "single"
	Podcast MetaType = "podcast"
	Music   MetaType = "music"
)

type Artist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Track struct {
	Id       int
	Title    string
	Duration time.Duration
	Artists  []Artist
}

type YaError struct {
	Err struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	} `json:"error"`
}

func (p *Playlist) UnmarshalJSON(data []byte) error {
	var tmp playlistResponse

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	p.Id = tmp.Id
	p.Kind = tmp.Kind
	p.Title = tmp.Title
	p.Duration = time.Millisecond * time.Duration(tmp.DurationMs)

	p.Tracks = make([]Track, 0, len(tmp.Tracks))
	for _, track := range tmp.Tracks {
		p.Tracks = append(p.Tracks, track.Track)
	}

	return nil
}

func (a *Album) UnmarshalJSON(data []byte) error {
	var tmp albumResponse

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	if tmp.Err != nil {
		*(a.Err) = *tmp.Err
		return nil
	}

	a.Id = tmp.Id
	a.Title = tmp.Title
	a.MetaType = tmp.MetaType
	a.Genre = tmp.Genre
	a.Artists = tmp.Artists

	a.Tracks = make([]Track, 0, tmp.TrackCount)
	for _, vol := range tmp.Volumes {
		a.Tracks = append(a.Tracks, vol...)
	}

	return nil
}

func (t *Track) UnmarshalJSON(data []byte) error {
	var tmp trackResponse

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	if tmp.Type != "music" {
		return ErrTrackIsNotMusic
	}

	t.Id = tmp.Id
	t.Title = tmp.Title
	t.Duration = time.Millisecond * time.Duration(tmp.DurationMs)
	t.Artists = tmp.Artists

	return nil
}

type playlistResponse struct {
	Id         int         `json:"uid"`
	Kind       int         `json:"kind"`
	Title      string      `json:"title"`
	DurationMs int         `json:"durationMs"`
	Tracks     []trackItem `json:"Tracks"`
}

type albumResponse struct {
	Id         int       `json:"id"`
	Err        *string   `json:"error"`
	Title      string    `json:"title"`
	MetaType   MetaType  `json:"metaType"`
	Genre      string    `json:"genre"`
	Artists    []Artist  `json:"artists"`
	TrackCount int       `json:"trackCount"`
	Volumes    [][]Track `json:"volumes"`
}

type trackResponse struct {
	Id         int      `json:"id"`
	Title      string   `json:"title"`
	DurationMs int      `json:"durationMs"`
	Artists    []Artist `json:"artists"`
	Type       string   `json:"type"`
}

type trackItem struct {
	Id    int   `json:"id"`
	Track Track `json:"track"`
}

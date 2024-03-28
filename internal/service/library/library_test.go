package library

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrackRegExp(t *testing.T) {
	type expected struct {
		trackId string
		albumId string
	}

	testCases := []struct {
		desc     string
		url      string
		expected expected
	}{
		{
			desc: "",
			url:  "https://music.yandex.ru/album/4867528/track/38184060",
			expected: expected{
				trackId: "38184060",
				albumId: "4867528",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			trackId, albumId := exctractTrackInfo(tC.url)
			assert.Equal(t, tC.expected.trackId, trackId)
			assert.Equal(t, tC.expected.albumId, albumId)
		})
	}
}

func TestAlbumRegExp(t *testing.T) {
	type expected struct {
		albumId string
	}

	testCases := []struct {
		desc     string
		url      string
		expected expected
	}{
		{
			desc: "",
			url:  "https://music.yandex.ru/album/6569250",
			expected: expected{
				albumId: "6569250",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			albumId := extractAlbumInfo(tC.url)
			assert.Equal(t, tC.expected.albumId, albumId)
		})
	}
}

func TestPlaylistRegExp(t *testing.T) {
	type expected struct {
		kind string
		user string
	}

	testCases := []struct {
		desc     string
		url      string
		expected expected
	}{
		{
			desc: "",
			url:  "https://music.yandex.com/users/remizova.av@phystech.edu/playlists/1007",
			expected: expected{
				kind: "1007",
				user: "remizova.av@phystech.edu",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			kind, user := exctractPlaylistInfo(tC.url)
			assert.Equal(t, tC.expected.user, user)
			assert.Equal(t, tC.expected.kind, kind)
		})
	}
}

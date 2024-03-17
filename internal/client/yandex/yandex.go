package client

import (
	"context"

	"github.com/GintGld/fizteh-radio-bot/internal/models"
)

type Client struct {
	token string
}

func New(
	token string,
) *Client {
	return &Client{
		token: token,
	}
}

func (c *Client) DownloadTrack(ctx context.Context, id int) (models.Media, error) {
	panic("not implemented")
}

func (c *Client) DownloadPlaylist(ctx context.Context, user string, id int) (models.Playlist, error) {
	panic("not implemented")
}

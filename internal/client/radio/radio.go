package client

import (
	"context"

	"github.com/GintGld/fizteh-radio-bot/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) GetToken(ctx context.Context, user models.User) (jwt.Token, error) {
	panic("not implemented")
}

func (c *Client) Search(ctx context.Context, token jwt.Token, filter models.MediaFilter) ([]models.Media, error) {
	panic("not implemented")
}

func (c *Client) NewMedia(ctx context.Context, token jwt.Token, media models.Media) error {
	panic("not implemented")
}

func (c *Client) NewTag(ctx context.Context, token jwt.Token, tag models.Tag) error {
	panic("not implemented")
}

func (c *Client) NewSegment(ctx context.Context, token jwt.Token, segm models.Segment) error {
	panic("not implemented")
}

func (c *Client) GetSchedule(ctx context.Context, token jwt.Token) ([]models.Segment, error) {
	panic("not implemented")
}

func (c *Client) GetConfig(ctx context.Context, token jwt.Token) (models.AutoDJConfig, error) {
	panic("not implemented")
}

func (c *Client) SetConfig(ctx context.Context, token jwt.Token, conf models.AutoDJConfig) error {
	panic("not implemented")
}

func (c *Client) StartAutoDJ(ctx context.Context, token jwt.Token) error {
	panic("not implemented")
}

func (c *Client) StopAutoDJ(ctx context.Context, token jwt.Token) error {
	panic("not implemented")
}

func (c *Client) IsAutoDJPlaying(ctx context.Context, token jwt.Token) (bool, error) {
	panic("not implemented")
}

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/GintGld/fizteh-radio-bot/internal/client"
	"github.com/GintGld/fizteh-radio-bot/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// TODO handle error messages.

type Client struct {
	addr      string
	c         *http.Client
	jwtParser *jwt.Parser
}

type HTTPError struct {
	Err string `json:"error"`
}

func New(
	addr string,
) *Client {
	return &Client{
		addr:      addr,
		c:         http.DefaultClient,
		jwtParser: new(jwt.Parser),
	}
}

func (c *Client) GetToken(ctx context.Context, user models.User) (jwt.Token, error) {
	const op = "Client.GetToken"

	url := fmt.Sprintf("%s/login", c.addr)

	bodyReq, err := json.Marshal(map[string]string{
		"login": user.Login,
		"pass":  user.Pass,
	})
	if err != nil {
		return jwt.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyReq))
	if err != nil {
		return jwt.Token{}, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return jwt.Token{}, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return jwt.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	switch resp.StatusCode {
	case 200:
		var form struct {
			Token string `json:"token"`
		}
		if err = json.Unmarshal(bodyResp, &form); err != nil {
			return jwt.Token{}, fmt.Errorf("%s: %w", op, err)
		}
		token, _, err := c.jwtParser.ParseUnverified(form.Token, jwt.MapClaims{})
		if err != nil {
			return jwt.Token{}, fmt.Errorf("%s: %w", op, err)
		}
		return *token, nil
	case 400:
		var e HTTPError
		if err := json.Unmarshal(bodyResp, &e); err != nil {
			return jwt.Token{}, fmt.Errorf("%s: %s", op, string(bodyResp))
		}
		return jwt.Token{}, fmt.Errorf("%s: returned error %s", op, e.Err)
	case 500:
		return jwt.Token{}, client.ErrInternalServerError
	default:
		return jwt.Token{}, fmt.Errorf("%s: unknown return status %d", op, resp.StatusCode)
	}
}

func (c *Client) Search(ctx context.Context, token jwt.Token, filter models.MediaFilter) ([]models.Media, error) {
	const op = "Client.Search"

	url := fmt.Sprintf("%s/library/media", c.addr)

	query := make([]string, 0, 3)
	if filter.Name != "" {
		query = append(query, fmt.Sprintf("name=%s", filter.Name))
	}
	if filter.Author != "" {
		query = append(query, fmt.Sprintf("author=%s", filter.Author))
	}
	if len(filter.Tags) > 0 {
		query = append(query, fmt.Sprintf("tags=%s", strings.Join(filter.Tags, ",")))
	}

	if len(query) > 0 {
		url += "?" + strings.Join(query, "&")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return []models.Media{}, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "Bearer "+token.Raw)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return []models.Media{}, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return []models.Media{}, fmt.Errorf("%s: %w", op, err)
	}

	switch resp.StatusCode {
	case 200:
		var resp struct {
			Library []models.Media `json:"library"`
		}
		if err := json.Unmarshal(bodyResp, &resp); err != nil {
			return []models.Media{}, fmt.Errorf("%s: %w", op, err)
		}
		return resp.Library, nil
	case 400:
		var e HTTPError
		if err := json.Unmarshal(bodyResp, &e); err != nil {
			return []models.Media{}, fmt.Errorf("%s: %s", op, string(bodyResp))
		}
		return []models.Media{}, fmt.Errorf("%s: returned error %s", op, e.Err)
	case 401:
		return []models.Media{}, client.ErrNotAuthorized
	case 500:
		return []models.Media{}, client.ErrInternalServerError
	default:
		return []models.Media{}, fmt.Errorf("%s: unknown return status %d", op, resp.StatusCode)
	}
}

func (c *Client) NewMedia(ctx context.Context, token jwt.Token, media models.Media) error {
	const op = "Client.NewMedia"

	url := fmt.Sprintf("%s/library/media", c.addr)

	jsonBytes, err := json.Marshal(map[string]any{
		"name":   media.Name,
		"author": media.Author,
		"tags":   media.Tags,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	file, err := os.Open(media.SourcePath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("media", string(jsonBytes)); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	part, err := writer.CreateFormFile("source", media.SourcePath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = io.Copy(part, file); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "Bearer "+token.Raw)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.c.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	switch resp.StatusCode {
	case 200:
		var jsonResp struct {
			Id int64 `json:"id"`
		}
		if err := json.Unmarshal(bodyResp, &jsonResp); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	case 400:
		var e HTTPError
		if err := json.Unmarshal(bodyResp, &e); err != nil {
			return fmt.Errorf("%s: %s", op, string(bodyResp))
		}
		return fmt.Errorf("%s: returned error %s", op, e.Err)
	case 401:
		return client.ErrNotAuthorized
	case 500:
		return client.ErrInternalServerError
	default:
		return fmt.Errorf("%s: unknown return status %d", op, resp.StatusCode)
	}
}

func (c *Client) AllTags(ctx context.Context, token jwt.Token) (models.TagList, error) {
	const op = "Client.AllTags"

	url := fmt.Sprintf("%s/library/tag", c.addr)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return models.TagList{}, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "Bearer "+token.Raw)

	resp, err := c.c.Do(req)
	if err != nil {
		return models.TagList{}, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.TagList{}, fmt.Errorf("%s: %w", op, err)
	}

	switch resp.StatusCode {
	case 200:
		var resp struct {
			Tags models.TagList `json:"tags"`
		}
		if err := json.Unmarshal(bodyResp, &resp); err != nil {
			return models.TagList{}, fmt.Errorf("%s: %w", op, err)
		}
		return resp.Tags, nil
	case 400:
		var e HTTPError
		if err := json.Unmarshal(bodyResp, &e); err != nil {
			return models.TagList{}, fmt.Errorf("%s: %s", op, string(bodyResp))
		}
		return models.TagList{}, fmt.Errorf("%s: returned error %s", op, e.Err)
	case 401:
		return models.TagList{}, client.ErrNotAuthorized
	case 500:
		return models.TagList{}, client.ErrInternalServerError
	default:
		return models.TagList{}, fmt.Errorf("%s: unknown return status %d", op, resp.StatusCode)
	}
}

func (c *Client) NewTag(ctx context.Context, token jwt.Token, tag models.Tag) (int64, error) {
	const op = "Client.NewTag"

	url := fmt.Sprintf("%s/library/tag", c.addr)

	bodyReq, err := json.Marshal(map[string]any{
		"tag": tag,
	})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyReq))
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "Bearer "+token.Raw)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	switch resp.StatusCode {
	case 200:
		var resp struct {
			Id int64 `json:"id"`
		}
		if err := json.Unmarshal(bodyResp, &resp); err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		return resp.Id, nil
	case 400:
		var e HTTPError
		if err := json.Unmarshal(bodyResp, &e); err != nil {
			return 0, fmt.Errorf("%s: %s", op, string(bodyResp))
		}
		return 0, fmt.Errorf("%s: returned error %s", op, e.Err)
	case 401:
		return 0, client.ErrNotAuthorized
	case 500:
		return 0, client.ErrInternalServerError
	default:
		return 0, fmt.Errorf("%s: unknown return status %d", op, resp.StatusCode)
	}
}

func (c *Client) NewSegment(ctx context.Context, token jwt.Token, segm models.Segment) error {
	const op = "Client.NewSegment"

	url := fmt.Sprintf("%s/schedule", c.addr)

	bodyReq, err := json.Marshal(map[string]any{
		"segment": map[string]any{
			"mediaID":   segm.Media.ID,
			"Start":     segm.Start,
			"BeginCut":  segm.BeginCut,
			"StopCut":   segm.StopCut,
			"Protected": segm.Protected,
		},
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyReq))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "Bearer "+token.Raw)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	switch resp.StatusCode {
	case 200:
		var resp struct {
			Id int64 `json:"id"`
		}
		if err := json.Unmarshal(bodyResp, &resp); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	case 400:
		var e HTTPError
		if err := json.Unmarshal(bodyResp, &e); err != nil {
			return fmt.Errorf("%s: %s", op, string(bodyResp))
		}
		return fmt.Errorf("%s: returned error %s", op, e.Err)
	case 401:
		return client.ErrNotAuthorized
	case 500:
		return client.ErrInternalServerError
	default:
		return fmt.Errorf("%s: unknown return status %d", op, resp.StatusCode)
	}
}

func (c *Client) GetSchedule(ctx context.Context, token jwt.Token) ([]models.Segment, error) {
	const op = "Client.GetSchedule"

	url := fmt.Sprintf("%s/schedule?start=%d", c.addr, time.Now().Unix())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return []models.Segment{}, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "Bearer "+token.Raw)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return []models.Segment{}, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return []models.Segment{}, fmt.Errorf("%s: %w", op, err)
	}

	switch resp.StatusCode {
	case 200:
		var resp struct {
			Segments []models.Segment `json:"segments"`
		}
		if err := json.Unmarshal(bodyResp, &resp); err != nil {
			return []models.Segment{}, fmt.Errorf("%s: %w", op, err)
		}
		return resp.Segments, nil
	case 400:
		var e HTTPError
		if err := json.Unmarshal(bodyResp, &e); err != nil {
			return []models.Segment{}, fmt.Errorf("%s: %s", op, string(bodyResp))
		}
		return []models.Segment{}, fmt.Errorf("%s: returned error %s", op, e.Err)
	case 401:
		return []models.Segment{}, client.ErrNotAuthorized
	case 500:
		return []models.Segment{}, client.ErrInternalServerError
	default:
		return []models.Segment{}, fmt.Errorf("%s: unknown return status %d", op, resp.StatusCode)
	}
}

func (c *Client) GetConfig(ctx context.Context, token jwt.Token) (models.AutoDJConfig, error) {
	const op = "Client.GetConfig"

	url := fmt.Sprintf("%s/dj/config", c.addr)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return models.AutoDJConfig{}, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "Bearer "+token.Raw)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return models.AutoDJConfig{}, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.AutoDJConfig{}, fmt.Errorf("%s: %w", op, err)
	}

	switch resp.StatusCode {
	case 200:
		var resp struct {
			Config models.AutoDJConfig `json:"config"`
		}
		if err := json.Unmarshal(bodyResp, &resp); err != nil {
			return models.AutoDJConfig{}, fmt.Errorf("%s: %w", op, err)
		}
		return resp.Config, nil
	case 401:
		return models.AutoDJConfig{}, client.ErrNotAuthorized
	case 500:
		return models.AutoDJConfig{}, client.ErrInternalServerError
	default:
		return models.AutoDJConfig{}, fmt.Errorf("%s: unknown return status %d", op, resp.StatusCode)
	}
}

func (c *Client) SetConfig(ctx context.Context, token jwt.Token, conf models.AutoDJConfig) error {
	const op = "Client.GetConfig"

	url := fmt.Sprintf("%s/dj/config", c.addr)

	bodyReq, err := json.Marshal(map[string]any{
		"config": conf,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyReq))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "Bearer "+token.Raw)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return nil
	case 401:
		return client.ErrNotAuthorized
	case 500:
		return client.ErrInternalServerError
	default:
		return fmt.Errorf("%s: unknown return status %d", op, resp.StatusCode)
	}
}

func (c *Client) StartAutoDJ(ctx context.Context, token jwt.Token) error {
	const op = "Client.NewSegment"

	url := fmt.Sprintf("%s/dj/start", c.addr)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "Bearer "+token.Raw)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return nil
	case 401:
		return client.ErrNotAuthorized
	case 500:
		return client.ErrInternalServerError
	default:
		return fmt.Errorf("%s: unknown return status %d", op, resp.StatusCode)
	}
}

func (c *Client) StopAutoDJ(ctx context.Context, token jwt.Token) error {
	const op = "Client.NewSegment"

	url := fmt.Sprintf("%s/dj/stop", c.addr)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "Bearer "+token.Raw)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return nil
	case 401:
		return client.ErrNotAuthorized
	case 500:
		return client.ErrInternalServerError
	default:
		return fmt.Errorf("%s: unknown return status %d", op, resp.StatusCode)
	}
}

func (c *Client) IsAutoDJPlaying(ctx context.Context, token jwt.Token) (bool, error) {
	const op = "Client.NewSegment"

	url := fmt.Sprintf("%s/dj/status", c.addr)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "Bearer "+token.Raw)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	switch resp.StatusCode {
	case 200:
		var resp struct {
			IsPlaying bool `json:"isPlaying"`
		}
		if err := json.Unmarshal(bodyResp, &resp); err != nil {
			return false, fmt.Errorf("%s: %w", op, err)
		}
		return resp.IsPlaying, nil
	case 401:
		return false, client.ErrNotAuthorized
	case 500:
		return false, client.ErrInternalServerError
	default:
		return false, fmt.Errorf("%s: unknown return status %d", op, resp.StatusCode)
	}
}

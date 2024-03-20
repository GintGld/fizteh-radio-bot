package client

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/GintGld/fizteh-radio-bot/internal/client"
	yamodels "github.com/GintGld/fizteh-radio-bot/internal/models/yandex"
)

const (
	proxyAddr = "api.music.yandex.net"
)

type Client struct {
	c      *http.Client
	token  string
	tmpDir string
}

func New(
	token string,
	tmpDir string,
) *Client {
	return &Client{
		c:      http.DefaultClient,
		token:  token,
		tmpDir: tmpDir,
	}
}

// Album returns album info
func (c *Client) Album(ctx context.Context, id string) (yamodels.Album, error) {
	const op = "Client.TrackInfo"

	url := fmt.Sprintf("https://%s/albums/%s/with-tracks", proxyAddr, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return yamodels.Album{}, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "OAuth "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return yamodels.Album{}, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return yamodels.Album{}, fmt.Errorf("%s: %w", op, err)
	}

	switch resp.StatusCode {
	case 200:
		var res struct {
			Res yamodels.Album `json:"result"`
		}
		if err := json.Unmarshal(bodyResp, &res); err != nil {
			return yamodels.Album{}, fmt.Errorf("%s: %w", op, err)
		}
		return res.Res, nil
	case 400:
		var res yamodels.YaError
		if err := json.Unmarshal(bodyResp, &res); err != nil {
			fmt.Printf("%s\n", err.Error())
			return yamodels.Album{}, fmt.Errorf("%s: status 400. Resp body: %s", op, string(bodyResp))
		}
		return yamodels.Album{}, fmt.Errorf("%s: status 400. Error result %+v", op, res)
	case 401:
		return yamodels.Album{}, client.ErrNotAuthorized
	case 500:
		return yamodels.Album{}, client.ErrInternalServerError
	default:
		return yamodels.Album{}, fmt.Errorf("%s: unknown return status %d. Resp body: %s", op, resp.StatusCode, string(bodyResp))
	}
}

// Playlist returns playlist info
func (c *Client) Playlist(ctx context.Context, user string, id string) (yamodels.Playlist, error) {
	const op = "Client.Playlist"

	url := fmt.Sprintf("https://%s/users/%s/playlists/%s", proxyAddr, user, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return yamodels.Playlist{}, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "OAuth "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return yamodels.Playlist{}, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return yamodels.Playlist{}, fmt.Errorf("%s: %w", op, err)
	}

	switch resp.StatusCode {
	case 200:
		var resp struct {
			Res yamodels.Playlist `json:"result"`
		}
		if err := json.Unmarshal(bodyResp, &resp); err != nil {
			return yamodels.Playlist{}, fmt.Errorf("%s: %w", op, err)
		}
		return resp.Res, nil
	case 400:
		var res yamodels.YaError
		if err := json.Unmarshal(bodyResp, &res); err != nil {
			fmt.Printf("%s\n", err.Error())
			return yamodels.Playlist{}, fmt.Errorf("%s: status 400. Resp body: %s", op, string(bodyResp))
		}
		return yamodels.Playlist{}, fmt.Errorf("%s: status 400. Error result %+v", op, res)
	case 401:
		return yamodels.Playlist{}, client.ErrNotAuthorized
	case 500:
		return yamodels.Playlist{}, client.ErrInternalServerError
	default:
		return yamodels.Playlist{}, fmt.Errorf("%s: unknown return status %d. Resp body: %s", op, resp.StatusCode, string(bodyResp))
	}
}

func (c *Client) DownloadInfo(ctx context.Context, id string) ([]yamodels.DownloadInfo, error) {
	const op = "Client.DownloadTrack"

	url := fmt.Sprintf("https://%s/tracks/%s/download-info", proxyAddr, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return []yamodels.DownloadInfo{}, fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "OAuth "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return []yamodels.DownloadInfo{}, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return []yamodels.DownloadInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	switch resp.StatusCode {
	case 200:
		var resp struct {
			Res []yamodels.DownloadInfo `json:"result"`
		}
		if err := json.Unmarshal(bodyResp, &resp); err != nil {
			return []yamodels.DownloadInfo{}, fmt.Errorf("%s: %w", op, err)
		}
		return resp.Res, nil
	case 400:
		var res yamodels.YaError
		if err := json.Unmarshal(bodyResp, &res); err != nil {
			fmt.Printf("%s\n", err.Error())
			return []yamodels.DownloadInfo{}, fmt.Errorf("%s: status 400. Resp body: %s", op, string(bodyResp))
		}
		return []yamodels.DownloadInfo{}, fmt.Errorf("%s: status 400. Error result %+v", op, res)
	case 401:
		return []yamodels.DownloadInfo{}, client.ErrNotAuthorized
	case 500:
		return []yamodels.DownloadInfo{}, client.ErrInternalServerError
	default:
		return []yamodels.DownloadInfo{}, fmt.Errorf("%s: unknown return status %d. Resp body: %s", op, resp.StatusCode, string(bodyResp))
	}
}

func (c *Client) DirectLink(ctx context.Context, url string) (string, error) {
	const op = "Client.getDirectLink"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "OAuth "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	switch resp.StatusCode {
	case 200:
		var resp yamodels.DownloadInfoXMLResponse
		if err := xml.Unmarshal(bodyResp, &resp); err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
		return resp.BuildLink(), nil
	case 400:
		var res yamodels.YaError
		if err := json.Unmarshal(bodyResp, &res); err != nil {
			fmt.Printf("%s\n", err.Error())
			return "", fmt.Errorf("%s: status 400. Resp body: %s", op, string(bodyResp))
		}
		return "", fmt.Errorf("%s: status 400. Error result %+v", op, res)
	case 401:
		return "", client.ErrNotAuthorized
	case 500:
		return "", client.ErrInternalServerError
	default:
		return "", fmt.Errorf("%s: unknown return status %d. Resp body: %s", op, resp.StatusCode, string(bodyResp))
	}
}

func (c *Client) DownloadTrack(ctx context.Context, url string) (string, error) {
	const op = "Client.downloadTrack"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	req.Header.Set("Authorization", "OAuth "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.c.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	out, err := os.CreateTemp(c.tmpDir, "ya-track-*.mp3")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return out.Name(), nil
}

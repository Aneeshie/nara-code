package firecracker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	base       string
}

type LoggerConfig struct {
	/// Named pipe or file used as output for logs.
	LogPath string `json:"log_path,omitempty"`
	/// The level of the Logger.
	Level LevelFilter `json:"level,omitempty"`
	/// Whether to show the log level in the log.
	ShowLevel bool `json:"show_level,omitempty"`
	/// Whether to show the log origin in the log.
	ShowLogOrigin bool `json:"show_log_origin,omitempty"`
}

type LevelFilter string

const (
	LevelFilterOff   LevelFilter = "Off"
	LevelFilterTrace LevelFilter = "Trace"
	LevelFilterDebug LevelFilter = "Debug"
	LevelFilterInfo  LevelFilter = "Info"
	LevelFilterWarn  LevelFilter = "Warn"
	LevelFilterError LevelFilter = "Error"
)

func NewClient(socketPath string) *Client {
	transport := http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", socketPath)
		},
	}

	client := http.Client{
		Transport: &transport,
	}

	return &Client{
		httpClient: &client,
		base:       "http://localhost/",
	}

}

func (c *Client) SetLogger(cfg *LoggerConfig) error {
	//marshal the config first
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", c.base+"logger", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil

}

// func (c *Client) SetBootSource() error
// func (c *Client) SetRootFs() error
// func (c *Client) SetActions() error

package betting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

const (
	tokenCreatePath = "/token/create"
)

type Client struct {
	url        string
	httpClient *http.Client
	log        *zap.Logger
}

func NewClient(u string, httpClient *http.Client, logger *zap.Logger) *Client {
	return &Client{
		url:        u,
		httpClient: httpClient,
		log:        logger,
	}
}

func (c *Client) GetToken(ctx context.Context, data map[string]any) (string, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("marshal request body: %w", err)
	}

	destinationURL, err := url.JoinPath(c.url, tokenCreatePath)
	if err != nil {
		return "", fmt.Errorf("generate destination URL: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		destinationURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}

	defer response.Body.Close()

	rawBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		c.log.Error("failed to create token", zap.ByteString("response_body", rawBody))

		return "", fmt.Errorf("invalid response status code: %d", response.StatusCode)
	}

	res := struct {
		Token string `json:"token"`
	}{}

	if err := json.Unmarshal(rawBody, &res); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return res.Token, nil
}

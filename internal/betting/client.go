package betting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	destinationURL, err := url.JoinPath(c.url, tokenCreatePath)
	if err != nil {
		return "", fmt.Errorf("failed to generate destination URL: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		destinationURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid response status code: %d", response.StatusCode)
	}

	defer response.Body.Close()

	res := struct {
		Token string `json:"token"`
	}{}

	if err := json.NewDecoder(response.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return res.Token, nil
}

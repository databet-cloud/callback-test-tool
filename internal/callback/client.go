package callback

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"

	"go.uber.org/zap"
)

type RequestType string

const (
	betPlacePath                 = "/bet/place"
	betAcceptPath                = "/bet/accept"
	betDeclinePath               = "/bet/decline"
	betSettlePath                = "/bet/settle"
	betUnSettlePath              = "/bet/unsettle"
	betCashOutOrdersAcceptedPath = "/bet/cash-out-orders/accepted"
	betCashOutOrdersDeclinedPath = "/bet/cash-out-orders/declined"

	BetPlaceRequestType                 RequestType = "place"
	BetAcceptRequestType                RequestType = "accept"
	BetDeclineRequestType               RequestType = "decline"
	BetSettleRequestType                RequestType = "settle"
	BetUnSettleRequestType              RequestType = "unsettle"
	BetCashOutOrdersAcceptedRequestType RequestType = "cash-out_accepted"
	BetCashOutOrdersDeclinedRequestType RequestType = "cash-out_declined"
)

type Client struct {
	url           string
	foreignParams map[string]any
	httpClient    *http.Client
	log           *zap.Logger
}

func NewClient(u string, foreignParams map[string]any, client *http.Client, log *zap.Logger) *Client {
	return &Client{
		url:           u,
		foreignParams: foreignParams,
		httpClient:    client,
		log:           log,
	}
}

func (c *Client) SendCallback(ctx context.Context, data *Data) (*http.Response, error) {
	switch data.RequestType {
	case BetPlaceRequestType:
		return c.sendRequest(ctx, betPlacePath, data)
	case BetAcceptRequestType:
		return c.sendRequest(ctx, betAcceptPath, data)
	case BetDeclineRequestType:
		return c.sendRequest(ctx, betDeclinePath, data)
	case BetSettleRequestType:
		return c.sendRequest(ctx, betSettlePath, data)
	case BetUnSettleRequestType:
		return c.sendRequest(ctx, betUnSettlePath, data)
	case BetCashOutOrdersAcceptedRequestType:
		return c.sendRequest(ctx, betCashOutOrdersAcceptedPath, data)
	case BetCashOutOrdersDeclinedRequestType:
		return c.sendRequest(ctx, betCashOutOrdersDeclinedPath, data)
	}

	return nil, errors.New("invalid request type")
}

func (c *Client) sendRequest(ctx context.Context, path string, body any) (*http.Response, error) {
	destinationURL, err := url.JoinPath(c.url, path)
	if err != nil {
		return nil, fmt.Errorf("failed to build destination url: %w", err)
	}

	foreignParamsValue, err := json.Marshal(c.foreignParams)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal foreign params: %w", err)
	}

	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, destinationURL, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Foreign-Params", string(foreignParamsValue))
	req.Header.Set("Content-Type", "application/json")

	dumpRequest, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, fmt.Errorf("failed to dump request: %w", err)
	}

	c.log.Debug("Send callback", zap.String("request", string(dumpRequest)))

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if response.StatusCode != http.StatusNoContent {
		rawBody, _ := io.ReadAll(response.Body)

		return nil, fmt.Errorf("unknown status code %d, body %s", response.StatusCode, rawBody)
	}

	return response, nil
}

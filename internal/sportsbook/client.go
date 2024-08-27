package sportsbook

import (
	"context"
	_ "embed"
	"fmt"

	"go.uber.org/zap"

	"github.com/machinebox/graphql"
)

var (
	//go:embed sportEventListByFilters.gql
	sportEventListByFiltersBody []byte
)

type Client struct {
	gqlClient *graphql.Client
	token     string
	log       *zap.Logger
}

func NewSportsBookClient(
	client *graphql.Client,
	token string,
	log *zap.Logger,
) *Client {
	return &Client{
		gqlClient: client,
		token:     token,
		log:       log,
	}
}

func (c *Client) SportEventsByFilter(ctx context.Context, offset, limit int) ([]SportEvent, error) {
	reg := graphql.NewRequest(string(sportEventListByFiltersBody))

	reg.Var("offset", offset)
	reg.Var("limit", limit)

	result := SportEventListByFilters{}

	if err := c.send(ctx, reg, &result); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return result.SportEventListByFilters.SportEvents, nil
}

func (c *Client) Token() string {
	return c.token
}

func (c *Client) send(ctx context.Context, req *graphql.Request, resp any) error {
	req.Header.Add("X-Auth-Token", c.token)

	return c.gqlClient.Run(ctx, req, resp)
}

package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/machinebox/graphql"
	"go.uber.org/zap"

	"github.com/databet-cloud/callback-test-tool/cmd/console/config"
	"github.com/databet-cloud/callback-test-tool/internal/betting"
	"github.com/databet-cloud/callback-test-tool/internal/sportsbook"
)

func MustCreateBettingClient(cfg config.Configuration, logger *zap.Logger) *betting.Client {
	httpClient, err := makeHttpClientWithTLSCertificate(cfg.Betting.Certificate)
	if err != nil {
		logger.Panic("Failed to create betting client", zap.Error(err))

		return nil
	}

	return betting.NewClient(cfg.Betting.URL, httpClient, logger)
}

func MustCreateSportsBookClient(cfg config.Configuration, token string, logger *zap.Logger) *sportsbook.Client {
	return sportsbook.NewSportsBookClient(
		graphql.NewClient(cfg.DataBetGQLURL, func(client *graphql.Client) {
			client.Log = func(s string) {
				logger.Sugar().Debug(s)
			}
		}),
		token,
		logger,
	)
}

func MustCreateLogger(cfg config.Configuration) *zap.Logger {
	if cfg.Debug {
		return zap.Must(zap.NewDevelopment())
	}

	return zap.Must(zap.NewProduction())
}

func makeHttpClientWithTLSCertificate(cert config.Certificate) (*http.Client, error) {
	certificate, err := tls.LoadX509KeyPair(cert.Path, cert.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %s", err)
	}

	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{Certificates: []tls.Certificate{certificate}}},
	}

	return client, nil
}

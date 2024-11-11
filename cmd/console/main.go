package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/databet-cloud/callback-test-tool/cmd/console/command"
	"github.com/databet-cloud/callback-test-tool/cmd/console/config"
	"github.com/databet-cloud/callback-test-tool/internal/balance"
	"github.com/databet-cloud/callback-test-tool/internal/calculator"
	"github.com/databet-cloud/callback-test-tool/internal/calculator/former"
	"github.com/databet-cloud/callback-test-tool/internal/callback"
	"github.com/databet-cloud/callback-test-tool/internal/prompt"
	"github.com/databet-cloud/callback-test-tool/internal/service"
)

func main() {
	var (
		ctx     = context.Background()
		rootCmd = &cobra.Command{
			Use:   "callback-test-tool",
			Short: "CLI tool to test callback server",
		}
		cfg = config.LoadConfig(rootCmd)
	)

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		run(ctx, cfg)
	}

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

//go:embed token_create_request.json
var rawTokenCreateRequest []byte

func run(ctx context.Context, cfg config.Configuration) {
	var (
		tokenCreateReq = map[string]any{}
		log            = MustCreateLogger(cfg)
		playerBalance  = balance.NewService(log)
		bettingClient  = MustCreateBettingClient(cfg, log.Named("betting"))
	)

	err := json.Unmarshal(rawTokenCreateRequest, &tokenCreateReq)
	if err != nil {
		panic(err)
	}

	authToken, err := bettingClient.GetToken(ctx, tokenCreateReq)
	if err != nil {
		log.Fatal("failed to get auth token", zap.Error(err))
	}

	log.Info("authenticated", zap.String("token", authToken))

	userSv := service.NewService(
		tokenCreateReq["player_id"].(string),
		playerBalance,
		MustCreateSportsBookClient(cfg, authToken, log.Named("sports_book")),
		callback.NewClient(cfg.CallbackServerURL, extractForeignParams(tokenCreateReq), http.DefaultClient, log),
		calculator.NewCalculator(calculator.NewRefundCalc(log, former.FormExpresses), log),
		log,
	)

	if err := playerBalance.DepositFloat(cfg.Balance); err != nil {
		log.Fatal("failed to deposit user balance", zap.Float64("amount", cfg.Balance), zap.Error(err))
	}

	prompt.ProcessCommands(command.Tree(ctx, userSv, cfg, log))
}

func extractForeignParams(tokenCreateReq map[string]any) map[string]any {
	v, ok := tokenCreateReq["params"]
	if !ok {
		return map[string]any{}
	}

	return v.(map[string]any)
}

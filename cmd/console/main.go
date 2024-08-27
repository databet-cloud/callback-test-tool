package main

import (
	"context"
	"net/http"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"gitlab.databet.one/b2b/callback-test-tool/cmd/console/command"
	"gitlab.databet.one/b2b/callback-test-tool/cmd/console/config"
	"gitlab.databet.one/b2b/callback-test-tool/internal/balance"
	"gitlab.databet.one/b2b/callback-test-tool/internal/calculator"
	"gitlab.databet.one/b2b/callback-test-tool/internal/calculator/former"
	"gitlab.databet.one/b2b/callback-test-tool/internal/callback"
	"gitlab.databet.one/b2b/callback-test-tool/internal/prompt"
	"gitlab.databet.one/b2b/callback-test-tool/internal/service"
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

func run(ctx context.Context, cfg config.Configuration) {
	var (
		foreignParams = map[string]any{
			"locale":    "en",
			"currency":  "usd",
			"player_id": cfg.PlayerID,
		}
		log           = MustCreateLogger(cfg)
		playerBalance = balance.NewService(log)
		bettingClient = MustCreateBettingClient(cfg, log.Named("betting"))
	)

	authToken, err := bettingClient.GetToken(ctx, foreignParams)
	if err != nil {
		log.Fatal("failed to get auth token", zap.Error(err))
	}

	log.Info("authenticated", zap.String("token", authToken))

	userSv := service.NewService(
		cfg.PlayerID,
		playerBalance,
		MustCreateSportsBookClient(cfg, authToken, log.Named("sports_book")),
		callback.NewClient(cfg.CallbackServerURL, foreignParams, http.DefaultClient, log),
		calculator.NewCalculator(calculator.NewRefundCalc(log, former.FormExpresses), log),
		log,
	)

	if err := playerBalance.DepositFloat(cfg.Balance); err != nil {
		log.Fatal("failed to deposit user balance", zap.Float64("amount", cfg.Balance), zap.Error(err))
	}

	prompt.ProcessCommands(command.Tree(ctx, userSv, cfg, log))
}

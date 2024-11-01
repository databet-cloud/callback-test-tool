package command

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"gitlab.databet.one/b2b/callback-test-tool/cmd/console/config"
	"gitlab.databet.one/b2b/callback-test-tool/internal/callback"
	"gitlab.databet.one/b2b/callback-test-tool/internal/prompt"
	"gitlab.databet.one/b2b/callback-test-tool/internal/service"
	"gitlab.databet.one/b2b/callback-test-tool/internal/storage"
)

const (
	selectBetLabel     = "Select bet (<id>:<state>_[<created>]:[<updated>])"
	selectCashOutLabel = "Select cash-out (<id>:<state>_[<created>]:[<updated>])"
)

func Tree(ctx context.Context, sv *service.Service, cfg config.Configuration, log *zap.Logger) *prompt.Tree {
	return &prompt.Tree{
		Label: "Select command",
		Commands: func() []*prompt.Command {
			return []*prompt.Command{
				player(sv, log),
				configCommand(cfg, log),
				placeBet(ctx, sv),
				acceptBet(ctx, sv),
				declineBet(ctx, sv),
				settleBet(ctx, sv, log),
				unSettleBet(ctx, sv),
				cashOut(ctx, sv),
				bets(sv),
				sentRequests(ctx, sv),
			}
		},
	}
}

func convert[T, V any](values []V, f func(V) T) []T {
	commands := make([]T, len(values))
	for i := range values {
		commands[i] = f(values[i])
	}

	return commands
}

func betDocLabel(doc *storage.Document[*callback.Data]) string {
	return fmt.Sprintf(
		"%s:%s_[%s]:[%s]",
		doc.Value.BetID,
		doc.Value.RequestType,
		doc.CreatedAt.Format(time.RFC3339),
		doc.UpdatedAt.Format(time.RFC3339),
	)
}

func cashOutDocLabel(doc *storage.Document[*callback.Data]) string {
	return fmt.Sprintf(
		"%s:%s:%s_[%s]:[%s]",
		doc.Value.BetID,
		doc.Value.CashOutOrderID,
		doc.Value.RequestType,
		doc.CreatedAt.Format(time.RFC3339),
		doc.UpdatedAt.Format(time.RFC3339),
	)
}

func printAsJSON(v any) {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		println(err.Error())
	}

	println(string(data))
}

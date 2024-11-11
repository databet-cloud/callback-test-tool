package command

import (
	"context"

	"github.com/databet-cloud/callback-test-tool/internal/callback"
	"github.com/databet-cloud/callback-test-tool/internal/prompt"
	"github.com/databet-cloud/callback-test-tool/internal/service"
	"github.com/databet-cloud/callback-test-tool/internal/storage"
)

func acceptBet(ctx context.Context, sv *service.Service) *prompt.Command {
	return &prompt.Command{
		Key: "accept bet",
		Tree: &prompt.Tree{
			Label:             selectBetLabel,
			ReturnAfterAction: true,
			Commands: func() []*prompt.Command {
				return convert(sv.Bets(callback.BetPlaceRequestType), func(d *storage.Document[*callback.Data]) *prompt.Command {
					return &prompt.Command{Key: betDocLabel(d), Action: func() { sv.AcceptBet(ctx, d.Value.BetID) }}
				})
			},
		},
	}
}

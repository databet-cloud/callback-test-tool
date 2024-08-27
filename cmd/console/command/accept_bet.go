package command

import (
	"context"

	"gitlab.databet.one/b2b/callback-test-tool/internal/callback"
	"gitlab.databet.one/b2b/callback-test-tool/internal/prompt"
	"gitlab.databet.one/b2b/callback-test-tool/internal/service"
	"gitlab.databet.one/b2b/callback-test-tool/internal/storage"
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

package command

import (
	"context"

	"gitlab.databet.one/b2b/callback-test-tool/internal/callback"
	"gitlab.databet.one/b2b/callback-test-tool/internal/prompt"
	"gitlab.databet.one/b2b/callback-test-tool/internal/service"
	"gitlab.databet.one/b2b/callback-test-tool/internal/storage"
)

func cashOut(ctx context.Context, sv *service.Service) *prompt.Command {
	return &prompt.Command{
		Key: "cash out",
		Tree: &prompt.Tree{
			Label: "Select cash out action",
			Commands: func() []*prompt.Command {
				return []*prompt.Command{
					{
						Key: "accept",
						Tree: &prompt.Tree{
							Label:             selectBetLabel,
							ReturnAfterAction: true,
							Commands: func() []*prompt.Command {
								bets := sv.Bets(callback.BetAcceptRequestType, callback.BetUnSettleRequestType)
								return convert(bets, func(d *storage.Document[*callback.Data]) *prompt.Command {
									return &prompt.Command{
										Key:    betDocLabel(d),
										Action: func() { sv.AcceptBetCashOut(ctx, d.Value.BetID) },
									}
								})
							},
						},
					},
					{
						Key: "decline",
						Tree: &prompt.Tree{
							Label:             selectBetLabel,
							ReturnAfterAction: true,
							Commands: func() []*prompt.Command {
								bets := sv.Bets(callback.BetCashOutOrdersAcceptedRequestType)
								return convert(bets, func(data *storage.Document[*callback.Data]) *prompt.Command {
									return &prompt.Command{
										Key:    betDocLabel(data),
										Action: func() { sv.DeclineBetCashOut(ctx, data.Value.BetID) },
									}
								})
							},
						},
					},
				}
			},
		},
	}
}

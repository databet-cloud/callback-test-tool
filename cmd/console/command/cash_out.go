package command

import (
	"context"
	"fmt"

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
							Label:             selectCashOutLabel,
							ReturnAfterAction: true,
							Commands: func() []*prompt.Command {
								cashOuts := sv.CashOuts(callback.BetCashOutOrdersAcceptedRequestType)
								return convert(cashOuts, func(data *storage.Document[*callback.Data]) *prompt.Command {
									return &prompt.Command{
										Key:    cashOutDocLabel(data),
										Action: func() { sv.DeclineBetCashOut(ctx, data.Value.BetID, data.Value.CashOutOrderID) },
									}
								})
							},
						},
					},
					{
						Key: "list",
						Tree: &prompt.Tree{
							Label: selectCashOutLabel,
							Commands: func() []*prompt.Command {
								return convert(sv.CashOuts(), func(d *storage.Document[*callback.Data]) *prompt.Command {
									return &prompt.Command{
										Key: cashOutDocLabel(d),
										Tree: &prompt.Tree{
											Label: fmt.Sprintf("cash-out %s", cashOutDocLabel(d)),
											Commands: func() []*prompt.Command {
												return []*prompt.Command{
													{
														Key:    "dump",
														Action: func() { printAsJSON(d.Value) },
													},
												}
											},
										},
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

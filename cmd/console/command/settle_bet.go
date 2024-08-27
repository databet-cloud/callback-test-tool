package command

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"gitlab.databet.one/b2b/callback-test-tool/internal/callback"
	"gitlab.databet.one/b2b/callback-test-tool/internal/prompt"
	"gitlab.databet.one/b2b/callback-test-tool/internal/service"
	"gitlab.databet.one/b2b/callback-test-tool/internal/sportsbook"
	"gitlab.databet.one/b2b/callback-test-tool/internal/storage"
)

func settleBet(ctx context.Context, sv *service.Service, log *zap.Logger) *prompt.Command {
	return &prompt.Command{
		Key: "settle bet",
		Tree: &prompt.Tree{
			Label:             selectBetLabel,
			ReturnAfterAction: true,
			Commands: func() []*prompt.Command {
				bets := sv.Bets(
					callback.BetAcceptRequestType,
					callback.BetUnSettleRequestType,
					callback.BetCashOutOrdersAcceptedRequestType,
					callback.BetCashOutOrdersDeclinedRequestType,
				)

				return convert(bets, func(d *storage.Document[*callback.Data]) *prompt.Command {
					return &prompt.Command{
						Key: betDocLabel(d),
						Action: func() {
							odds := make([]*callback.Odd, len(d.Value.PrivateOdds))
							for i, odd := range d.Value.PrivateOdds {
								label := fmt.Sprintf("Select status for [%s].[%s].[%s]", odd.MatchId, odd.MarketId, odd.OddId)
								status, err := prompt.Select(
									label,
									sportsbook.OddStatusWin,
									sportsbook.OddStatusHalfWin,
									sportsbook.OddStatusLoss,
									sportsbook.OddStatusHalfLoss,
									sportsbook.OddStatusRefunded,
									sportsbook.OddStatusRefundedManually,
								)
								if err != nil {
									log.Error("failed to set odd status", zap.Error(err))

									return
								}

								odds[i] = odd.WithStatus(status)
							}

							sv.SettleBet(ctx, d.Value.BetID, odds)
						},
					}
				})
			},
		},
	}
}

func unSettleBet(ctx context.Context, sv *service.Service) *prompt.Command {
	return &prompt.Command{
		Key: "unsettle bet",
		Tree: &prompt.Tree{
			Label:             selectBetLabel,
			ReturnAfterAction: true,
			Commands: func() []*prompt.Command {
				return convert(sv.Bets(callback.BetSettleRequestType), func(d *storage.Document[*callback.Data]) *prompt.Command {
					return &prompt.Command{Key: betDocLabel(d), Action: func() { sv.UnSettleBet(ctx, d.Value.BetID) }}
				})
			},
		},
	}
}

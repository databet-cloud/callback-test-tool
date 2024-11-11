package command

import (
	"context"
	"fmt"

	"github.com/databet-cloud/callback-test-tool/internal/callback"
	"github.com/databet-cloud/callback-test-tool/internal/prompt"
	"github.com/databet-cloud/callback-test-tool/internal/service"
	"github.com/databet-cloud/callback-test-tool/internal/storage"
)

func declineBet(ctx context.Context, sv *service.Service) *prompt.Command {
	return &prompt.Command{
		Key: "decline bet",
		Tree: &prompt.Tree{
			Label:             selectBetLabel,
			ReturnAfterAction: true,
			Commands: func() []*prompt.Command {
				requests := sv.Bets(callback.BetPlaceRequestType, callback.BetAcceptRequestType)

				return convert(requests, func(d *storage.Document[*callback.Data]) *prompt.Command {
					return &prompt.Command{
						Key: betDocLabel(d),
						Tree: &prompt.Tree{
							Label:             fmt.Sprintf("Select restriction to decline: %s", betDocLabel(d)),
							ReturnAfterAction: true,
							Commands: func() []*prompt.Command {
								rr := callback.GetAllBetRestrictions()

								return convert(rr, func(r callback.RestrictionType) *prompt.Command {
									return &prompt.Command{
										Key:    r.String(),
										Action: func() { sv.DeclineBet(ctx, d.Value.BetID, r) },
									}
								})
							},
						},
					}
				})
			},
		},
	}
}

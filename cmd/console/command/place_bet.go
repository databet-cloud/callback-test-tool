package command

import (
	"context"

	"github.com/databet-cloud/callback-test-tool/internal/callback"
	"github.com/databet-cloud/callback-test-tool/internal/prompt"
	"github.com/databet-cloud/callback-test-tool/internal/service"
)

func placeBet(ctx context.Context, sv *service.Service) *prompt.Command {
	return &prompt.Command{
		Key: "place bet",
		Tree: &prompt.Tree{
			Label: "Select bet type",
			Commands: func() []*prompt.Command {
				return []*prompt.Command{
					{
						Key:    "single",
						Action: func() { sv.PlaceBet(ctx, callback.SingleBetType, prompt.Float("Put bet amount")) },
					},
					{
						Key:    "express",
						Action: func() { sv.PlaceBet(ctx, callback.ExpressBetType, prompt.Float("Put bet amount")) },
					},
					{
						Key:    "system",
						Action: func() { sv.PlaceBet(ctx, callback.SystemBetType, prompt.Float("Put bet amount")) },
					},
				}
			},
			ReturnAfterAction: true,
		},
	}
}

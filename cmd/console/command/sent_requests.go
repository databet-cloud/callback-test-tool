package command

import (
	"context"
	"fmt"

	"gitlab.databet.one/b2b/callback-test-tool/internal/callback"
	"gitlab.databet.one/b2b/callback-test-tool/internal/prompt"
	"gitlab.databet.one/b2b/callback-test-tool/internal/service"
	"gitlab.databet.one/b2b/callback-test-tool/internal/storage"
)

func sentRequests(ctx context.Context, sv *service.Service) *prompt.Command {
	return &prompt.Command{
		Key: "sent requests",
		Tree: &prompt.Tree{
			Label:             "Select sent request(<bet>:<type>_[<created>]:[<updated>])",
			ReturnAfterAction: true,
			Commands: func() []*prompt.Command {
				return convert(sv.SentRequests(), func(d *storage.Document[*callback.Data]) *prompt.Command {
					return &prompt.Command{
						Key: betDocLabel(d),
						Tree: &prompt.Tree{
							Label: fmt.Sprintf("Request %s", betDocLabel(d)),
							Commands: func() []*prompt.Command {
								return []*prompt.Command{
									{
										Key:    "replay",
										Action: func() { sv.ReplayCallback(ctx, d.Value) },
									},
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
	}
}

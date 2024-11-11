package command

import (
	"fmt"

	"github.com/databet-cloud/callback-test-tool/internal/callback"
	"github.com/databet-cloud/callback-test-tool/internal/prompt"
	"github.com/databet-cloud/callback-test-tool/internal/service"
	"github.com/databet-cloud/callback-test-tool/internal/storage"
)

func bets(sv *service.Service) *prompt.Command {
	return &prompt.Command{
		Key: "bets",
		Tree: &prompt.Tree{
			Label: selectBetLabel,
			Commands: func() []*prompt.Command {
				return convert(sv.Bets(), func(d *storage.Document[*callback.Data]) *prompt.Command {
					return &prompt.Command{
						Key: betDocLabel(d),
						Tree: &prompt.Tree{
							Label: fmt.Sprintf("Bet %s", betDocLabel(d)),
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
	}
}

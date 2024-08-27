package command

import (
	"fmt"

	"gitlab.databet.one/b2b/callback-test-tool/internal/callback"
	"gitlab.databet.one/b2b/callback-test-tool/internal/prompt"
	"gitlab.databet.one/b2b/callback-test-tool/internal/service"
	"gitlab.databet.one/b2b/callback-test-tool/internal/storage"
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

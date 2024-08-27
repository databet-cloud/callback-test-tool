package command

import (
	"go.uber.org/zap"

	"gitlab.databet.one/b2b/callback-test-tool/internal/prompt"
	"gitlab.databet.one/b2b/callback-test-tool/internal/service"
)

func player(sv *service.Service, log *zap.Logger) *prompt.Command {
	return &prompt.Command{
		Key: "player",
		Action: func() {
			log.Info(
				"Player info",
				zap.String("id", sv.PlayerID()),
				zap.Any("balance", sv.PlayerBalance()),
				zap.String("token", sv.PlayerToken()),
			)
		},
	}
}

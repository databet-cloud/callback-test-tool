package command

import (
	"go.uber.org/zap"

	"gitlab.databet.one/b2b/callback-test-tool/cmd/console/config"
	"gitlab.databet.one/b2b/callback-test-tool/internal/prompt"
)

func configCommand(cfg config.Configuration, log *zap.Logger) *prompt.Command {
	return &prompt.Command{
		Key:    "config",
		Action: func() { log.Info("Configuration", zap.Any("cfg", cfg)) },
	}
}

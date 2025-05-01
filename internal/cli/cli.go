package cli

import (
	"fmt"

	"github.com/devusSs/minyls/internal/env"
	"github.com/devusSs/minyls/internal/log"
)

var e *env.Env

// TODO: add functions for storage & clip
func initialize() error {
	err := log.Setup()
	if err != nil {
		return fmt.Errorf("failed to setup log: %w", err)
	}

	log.Log().Debug().Str("func", "cli.initialize").Msg("setup log")

	e, err = env.Load()
	if err != nil {
		return fmt.Errorf("failed to load env: %w", err)
	}

	log.Log().Debug().Str("func", "cli.initialize").Any("env", e).Msg("loaded environment")

	return nil
}

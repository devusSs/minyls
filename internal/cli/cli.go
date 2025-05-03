package cli

import (
	"fmt"

	"github.com/devusSs/minyls/internal/clip"
	"github.com/devusSs/minyls/internal/env"
	"github.com/devusSs/minyls/internal/log"
	"github.com/devusSs/minyls/internal/storage"
)

var e *env.Env

// TODO: add functions for storage
func initialize() error {
	err := log.Setup()
	if err != nil {
		return fmt.Errorf("failed to setup log: %w", err)
	}

	log.Log().Debug().Str("func", "cli.initialize").Msg("setup log")

	err = clip.Init()
	if err != nil {
		return fmt.Errorf("failed to init clip: %w", err)
	}

	log.Log().Debug().Str("func", "cli.initialize").Msg("setup clip")

	e, err = env.Load()
	if err != nil {
		return fmt.Errorf("failed to load env: %w", err)
	}

	log.Log().Debug().Str("func", "cli.initialize").Any("env", e).Msg("loaded environment")

	err = storage.Init(e.MinioLinkExpiry)
	if err != nil {
		return fmt.Errorf("failed to init storage: %w", err)
	}

	log.Log().Debug().Str("func", "cli.Initialize").Msg("setup storage")

	return nil
}

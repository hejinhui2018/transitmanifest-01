package app

import (
	"errors"
	"fmt"
	"os"
)

func (a *App) Close() error {
	if a == nil || a.Store == nil {
		return nil
	}
	if err := a.Manifests.Snapshot(); err != nil {
		return fmt.Errorf("save application snapshot: %w", err)
	}
	return a.Store.Close()
}

func IsNotFound(err error) bool { return errors.Is(err, os.ErrNotExist) }

func (a *App) Snapshot() error { return a.Manifests.Snapshot() }

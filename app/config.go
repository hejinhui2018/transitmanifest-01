package app

import (
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	DataDir string
	Actor   string
}

func (c Config) Validate() error {
	if c.DataDir == "" {
		return errors.New("data directory is required")
	}
	if filepath.IsAbs(c.DataDir) == false {
		return errors.New("data directory must be absolute")
	}
	if c.Actor == "" {
		return errors.New("actor is required")
	}
	return nil
}

func FromEnvironment() Config {
	directory := os.Getenv("TRANSITMANIFEST_DATA_DIR")
	if directory == "" {
		directory = filepath.Join(".", "data")
	}
	actor := os.Getenv("TRANSITMANIFEST_ACTOR")
	if actor == "" {
		actor = "cli"
	}
	absolute, err := filepath.Abs(directory)
	if err == nil {
		directory = absolute
	}
	return Config{DataDir: directory, Actor: actor}
}

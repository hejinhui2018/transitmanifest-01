package app

import (
	"fmt"

	"transitmanifest/audit"
	"transitmanifest/exception"
	"transitmanifest/handoff"
	"transitmanifest/manifest"
	"transitmanifest/query"
	"transitmanifest/report"
	"transitmanifest/routing"
	"transitmanifest/scan"
	"transitmanifest/storage"
)

type App struct {
	Config     Config
	Store      *storage.Store
	Manifests  *manifest.Service
	Scans      *scan.Service
	Handoffs   *handoff.Service
	Routes     *routing.Service
	Exceptions *exception.Service
	Reports    *report.Service
	Audit      *audit.Service
	Query      *query.Service
}

func Open(config Config) (*App, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	store, err := storage.Open(config.DataDir)
	if err != nil {
		return nil, fmt.Errorf("open application storage: %w", err)
	}
	manifests, err := manifest.Open(store)
	if err != nil {
		_ = store.Close()
		return nil, fmt.Errorf("recover manifests: %w", err)
	}
	return &App{Config: config, Store: store, Manifests: manifests, Scans: scan.New(manifests), Handoffs: handoff.New(manifests), Routes: routing.New(manifests), Exceptions: exception.New(manifests), Reports: report.New(store), Audit: audit.New(store), Query: query.New(manifests)}, nil
}

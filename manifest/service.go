package manifest

import (
	"sync"
	"time"

	"transitmanifest/domain"
	"transitmanifest/storage"
)

type Service struct {
	mu    sync.RWMutex
	store *storage.Store
	now   func() time.Time
	state persistedState
}

func Open(store *storage.Store) (*Service, error) {
	state, err := recoverState(store)
	if err != nil {
		return nil, err
	}
	return &Service{store: store, now: time.Now, state: state}, nil
}

func (s *Service) Snapshot() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.store.Snapshot(s.state.clone())
}

func (s *Service) Get(id string) (domain.Manifest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.state.Manifests[id]
	if !ok {
		return domain.Manifest{}, domain.ErrManifestNotFound
	}
	return item.Clone(), nil
}

func (s *Service) List() []domain.Manifest {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return sortedManifests(s.state.Manifests)
}

func (s *Service) Scans() []domain.PackageScan {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.PackageScan, 0, len(s.state.Scans))
	for _, item := range s.state.Scans {
		result = append(result, item)
	}
	return result
}

func (s *Service) Exceptions() []domain.Exception {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return sortedExceptions(s.state.Exceptions)
}

func (s *Service) Routes() []domain.Route {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.Route, 0, len(s.state.Routes))
	for _, item := range s.state.Routes {
		result = append(result, item)
	}
	return result
}

func (s *Service) HasScan(scanID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.state.Scans[scanID]
	return ok
}

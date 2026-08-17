package storage

import (
	"fmt"
	"sync"

	"transitmanifest/domain"
)

// Store serializes event append and snapshot capture through a shared gate.
// Services can therefore capture a snapshot exactly at a known log position.
type Store struct {
	mu        sync.RWMutex
	directory string
	log       *Log
}

func Open(directory string) (*Store, error) {
	log, err := OpenLog(directory)
	if err != nil {
		return nil, err
	}
	return &Store{directory: directory, log: log}, nil
}

func (s *Store) Append(event domain.Event) (Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	record, err := s.log.Append(event)
	if err != nil {
		return Record{}, fmt.Errorf("persist event %s: %w", event.Type, err)
	}
	return record, nil
}

func (s *Store) Snapshot(payload any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sequence, hash := s.log.Position()
	if err := WriteSnapshot(s.directory, sequence, hash, payload); err != nil {
		return fmt.Errorf("persist snapshot at %d: %w", sequence, err)
	}
	return nil
}

func (s *Store) Recover(target any) (Snapshot, []Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	exists, err := SnapshotExists(s.directory)
	if err != nil {
		return Snapshot{}, nil, fmt.Errorf("inspect snapshot: %w", err)
	}
	var snapshot Snapshot
	if exists {
		snapshot, err = ReadSnapshot(s.directory, target)
		if err != nil {
			return Snapshot{}, nil, err
		}
	}
	records, err := ReadAfter(s.directory, snapshot.Sequence)
	if err != nil {
		return Snapshot{}, nil, err
	}
	if snapshot.Sequence > 0 {
		all, err := ReadAfter(s.directory, 0)
		if err != nil {
			return Snapshot{}, nil, err
		}
		if int(snapshot.Sequence) > len(all) || all[snapshot.Sequence-1].Checksum != snapshot.LogHash {
			return Snapshot{}, nil, fmt.Errorf("%w: log position mismatch", domain.ErrCorruptSnapshot)
		}
	}
	return snapshot, records, nil
}

func (s *Store) Records() ([]Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return Records(s.directory)
}

func (s *Store) Directory() string {
	return s.directory
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.log.Close()
}

package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Health struct {
	Directory       string    `json:"directory"`
	LogBytes        int64     `json:"log_bytes"`
	SnapshotBytes   int64     `json:"snapshot_bytes"`
	LogSequence     uint64    `json:"log_sequence"`
	LastChecksum    string    `json:"last_checksum"`
	SnapshotPresent bool      `json:"snapshot_present"`
	CheckedAt       time.Time `json:"checked_at"`
}

// HealthCheck verifies the complete append log and reports the durable
// position without exposing an open file descriptor to callers.
func (s *Store) HealthCheck() (Health, error) {
	result := Health{Directory: s.directory, CheckedAt: time.Now().UTC()}
	records, err := s.Records()
	if err != nil {
		return result, fmt.Errorf("verify storage health: %w", err)
	}
	if len(records) > 0 {
		last := records[len(records)-1]
		result.LogSequence = last.Sequence
		result.LastChecksum = last.Checksum
	}
	logInfo, err := os.Stat(filepath.Join(s.directory, eventLogName))
	if err != nil {
		return result, fmt.Errorf("stat event log: %w", err)
	}
	result.LogBytes = logInfo.Size()
	snapshotInfo, err := os.Stat(filepath.Join(s.directory, snapshotName))
	if err == nil {
		result.SnapshotPresent = true
		result.SnapshotBytes = snapshotInfo.Size()
	} else if !os.IsNotExist(err) {
		return result, fmt.Errorf("stat snapshot: %w", err)
	}
	return result, nil
}

func (s *Store) DurablePosition() (uint64, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.log.Position()
}

func (s *Store) EventCount() (int, error) {
	records, err := s.Records()
	if err != nil {
		return 0, err
	}
	return len(records), nil
}

func (s *Store) HasSnapshot() (bool, error) {
	return SnapshotExists(s.directory)
}

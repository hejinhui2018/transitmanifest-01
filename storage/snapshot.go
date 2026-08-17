package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"transitmanifest/domain"
)

const snapshotName = "snapshot.json"

type Snapshot struct {
	Sequence uint64          `json:"sequence"`
	LogHash  string          `json:"log_hash"`
	Payload  json.RawMessage `json:"payload"`
	Checksum string          `json:"checksum"`
}

func WriteSnapshot(directory string, sequence uint64, logHash string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal snapshot payload: %w", err)
	}
	snapshot := Snapshot{Sequence: sequence, LogHash: logHash, Payload: data}
	snapshot.Checksum = snapshotChecksum(snapshot)
	encoded, err := json.Marshal(snapshot)
	if err != nil {
		return fmt.Errorf("marshal snapshot envelope: %w", err)
	}
	encoded = append(encoded, '\n')
	if err := os.MkdirAll(directory, 0o750); err != nil {
		return fmt.Errorf("create snapshot directory: %w", err)
	}
	temporary, err := os.CreateTemp(directory, ".snapshot-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary snapshot: %w", err)
	}
	temporaryName := temporary.Name()
	committed := false
	defer func() {
		_ = temporary.Close()
		if !committed {
			_ = os.Remove(temporaryName)
		}
	}()
	if err := temporary.Chmod(0o640); err != nil {
		return fmt.Errorf("set snapshot permissions: %w", err)
	}
	if _, err := temporary.Write(encoded); err != nil {
		return fmt.Errorf("write temporary snapshot: %w", err)
	}
	if err := temporary.Sync(); err != nil {
		return fmt.Errorf("sync temporary snapshot: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary snapshot: %w", err)
	}
	target := filepath.Join(directory, snapshotName)
	if runtime.GOOS == "windows" {
		_ = os.Remove(target)
	}
	if err := os.Rename(temporaryName, target); err != nil {
		return fmt.Errorf("commit snapshot: %w", err)
	}
	committed = true
	if err := syncDirectory(directory); err != nil {
		return fmt.Errorf("sync snapshot directory: %w", err)
	}
	return nil
}

func ReadSnapshot(directory string, target any) (Snapshot, error) {
	path := filepath.Join(directory, snapshotName)
	file, err := os.Open(path)
	if err != nil {
		return Snapshot{}, err
	}
	defer file.Close()
	encoded, err := io.ReadAll(io.LimitReader(file, 64*1024*1024))
	if err != nil {
		return Snapshot{}, fmt.Errorf("read snapshot: %w", err)
	}
	var snapshot Snapshot
	if err := json.Unmarshal(encoded, &snapshot); err != nil {
		return Snapshot{}, fmt.Errorf("%w: decode envelope: %v", domain.ErrCorruptSnapshot, err)
	}
	if snapshot.Checksum != snapshotChecksum(snapshot) {
		return Snapshot{}, fmt.Errorf("%w: checksum mismatch", domain.ErrCorruptSnapshot)
	}
	if err := json.Unmarshal(snapshot.Payload, target); err != nil {
		return Snapshot{}, fmt.Errorf("%w: decode payload: %v", domain.ErrCorruptSnapshot, err)
	}
	return snapshot, nil
}

func SnapshotExists(directory string) (bool, error) {
	_, err := os.Stat(filepath.Join(directory, snapshotName))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func snapshotChecksum(snapshot Snapshot) string {
	digest := sha256.New()
	fmt.Fprintf(digest, "%d\n%s\n", snapshot.Sequence, snapshot.LogHash)
	digest.Write(snapshot.Payload)
	return hex.EncodeToString(digest.Sum(nil))
}

func syncDirectory(directory string) error {
	handle, err := os.Open(directory)
	if err != nil {
		return err
	}
	defer handle.Close()
	err = handle.Sync()
	if runtime.GOOS == "windows" && errors.Is(err, os.ErrInvalid) {
		return nil
	}
	return err
}

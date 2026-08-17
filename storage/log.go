package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"transitmanifest/domain"
)

const eventLogName = "events.log"

type Log struct {
	mu       sync.Mutex
	path     string
	file     *os.File
	sequence uint64
	lastHash string
	now      func() time.Time
}

func OpenLog(directory string) (*Log, error) {
	if err := os.MkdirAll(directory, 0o750); err != nil {
		return nil, fmt.Errorf("create storage directory: %w", err)
	}
	path := filepath.Join(directory, eventLogName)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o640)
	if err != nil {
		return nil, fmt.Errorf("open event log: %w", err)
	}
	log := &Log{path: path, file: file, now: time.Now}
	if err := log.inspect(); err != nil {
		_ = file.Close()
		return nil, err
	}
	return log, nil
}

func (l *Log) inspect() error {
	if _, err := l.file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("seek event log: %w", err)
	}
	scanner := bufio.NewScanner(l.file)
	buffer := make([]byte, 64*1024)
	scanner.Buffer(buffer, 4*1024*1024)
	var sequence uint64
	previous := ""
	for scanner.Scan() {
		var record Record
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return fmt.Errorf("%w: decode record %d: %v", domain.ErrCorruptLog, sequence+1, err)
		}
		if err := record.validate(sequence+1, previous); err != nil {
			return err
		}
		sequence = record.Sequence
		previous = record.Checksum
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan event log: %w", err)
	}
	l.sequence = sequence
	l.lastHash = previous
	_, err := l.file.Seek(0, io.SeekEnd)
	return err
}

func (l *Log) Append(event domain.Event) (Record, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file == nil {
		return Record{}, errors.New("event log is closed")
	}
	if err := domain.ValidateType(event.Type); err != nil {
		return Record{}, err
	}
	record := newRecord(l.sequence+1, l.lastHash, event, l.now())
	record.Checksum = record.calculateChecksum()
	encoded, err := json.Marshal(record)
	if err != nil {
		return Record{}, fmt.Errorf("encode event record: %w", err)
	}
	encoded = append(encoded, '\n')
	if _, err := l.file.Write(encoded); err != nil {
		return Record{}, fmt.Errorf("append event record: %w", err)
	}
	if err := l.file.Sync(); err != nil {
		return Record{}, fmt.Errorf("sync event record: %w", err)
	}
	l.sequence = record.Sequence
	l.lastHash = record.Checksum
	return record, nil
}

func (l *Log) Position() (uint64, string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.sequence, l.lastHash
}

func (l *Log) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}

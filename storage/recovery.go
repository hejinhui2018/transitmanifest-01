package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// ReadAfter validates the complete checksum chain before returning records.
// A caller cannot accidentally recover from an unverified suffix.
func ReadAfter(directory string, sequence uint64) ([]Record, error) {
	path := filepathJoin(directory, eventLogName)
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("open event log for recovery: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buffer := make([]byte, 64*1024)
	scanner.Buffer(buffer, 4*1024*1024)
	result := make([]Record, 0)
	var current uint64
	previous := ""
	for scanner.Scan() {
		var record Record
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, fmt.Errorf("decode recovery record %d: %w", current+1, err)
		}
		if err := record.validate(current+1, previous); err != nil {
			return nil, err
		}
		current = record.Sequence
		previous = record.Checksum
		if record.Sequence > sequence {
			result = append(result, record)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan recovery log: %w", err)
	}
	if sequence > current {
		return nil, fmt.Errorf("snapshot sequence %d exceeds log sequence %d", sequence, current)
	}
	return result, nil
}

// Records returns an independently allocated, verified view suitable for
// audits and reports. It intentionally does not expose the live file handle.
func Records(directory string) ([]Record, error) {
	return ReadAfter(directory, 0)
}

func filepathJoin(directory, name string) string {
	if len(directory) == 0 {
		return name
	}
	separator := string(os.PathSeparator)
	if directory[len(directory)-1:] == separator {
		return directory + name
	}
	return directory + separator + name
}

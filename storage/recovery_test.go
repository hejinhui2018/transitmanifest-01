package storage

import (
	"os"
	"testing"
	"time"

	"transitmanifest/domain"
)

func TestRecoveryRejectsTamperedLog(t *testing.T) {
	directory := t.TempDir()
	log, err := OpenLog(directory)
	if err != nil {
		t.Fatal(err)
	}
	event, _ := domain.NewEvent(domain.EventManifestClosed, "m-1", "tester", time.Now(), domain.ManifestClosedData{})
	if _, err := log.Append(event); err != nil {
		t.Fatal(err)
	}
	log.Close()
	path := directory + string(os.PathSeparator) + eventLogName
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for i := range data {
		if data[i] == 'm' {
			data[i] = 'n'
			break
		}
	}
	if err := os.WriteFile(path, data, 0o640); err != nil {
		t.Fatal(err)
	}
	if _, err := Records(directory); err == nil {
		t.Fatal("tampered log was accepted")
	}
}

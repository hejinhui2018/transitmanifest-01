package storage

import (
	"testing"
	"time"

	"transitmanifest/domain"
)

func TestLogAppendBuildsChecksumChain(t *testing.T) {
	directory := t.TempDir()
	log, err := OpenLog(directory)
	if err != nil {
		t.Fatal(err)
	}
	defer log.Close()
	event, err := domain.NewEvent(domain.EventManifestCreated, "m-1", "tester", time.Unix(1, 0), domain.ManifestCreated{TripID: "t", VehiclePlate: "p", Origin: "a", Destination: "b"})
	if err != nil {
		t.Fatal(err)
	}
	first, err := log.Append(event)
	if err != nil {
		t.Fatal(err)
	}
	second, err := log.Append(event)
	if err != nil {
		t.Fatal(err)
	}
	if first.Sequence != 1 || second.Sequence != 2 {
		t.Fatalf("unexpected sequence: %d %d", first.Sequence, second.Sequence)
	}
	if second.PreviousHash != first.Checksum || first.Checksum == "" || second.Checksum == "" {
		t.Fatal("checksum chain was not linked")
	}
}

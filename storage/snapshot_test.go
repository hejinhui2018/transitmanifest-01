package storage

import (
	"testing"
)

func TestSnapshotRoundTrip(t *testing.T) {
	directory := t.TempDir()
	if err := WriteSnapshot(directory, 7, "abc", map[string]any{"value": "kept"}); err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	snapshot, err := ReadSnapshot(directory, &payload)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.Sequence != 7 || snapshot.LogHash != "abc" || payload["value"] != "kept" {
		t.Fatalf("snapshot mismatch: %+v %+v", snapshot, payload)
	}
	exists, err := SnapshotExists(directory)
	if err != nil || !exists {
		t.Fatalf("snapshot missing: %v %v", exists, err)
	}
}

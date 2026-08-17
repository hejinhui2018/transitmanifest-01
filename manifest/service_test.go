package manifest

import (
	"testing"
	"time"

	"transitmanifest/storage"
)

func TestManifestLifecyclePersists(t *testing.T) {
	store, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	service, err := Open(store)
	if err != nil {
		t.Fatal(err)
	}
	created, err := service.Create("m-1", "trip-1", "plate-1", "hub-a", "hub-b", "tester", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if created.Status != "open" {
		t.Fatalf("unexpected status %s", created.Status)
	}
	closed, err := service.Close("m-1", "loaded", "tester", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if closed.Status != "closed" || closed.ClosedAt == nil {
		t.Fatalf("manifest was not closed: %+v", closed)
	}
	if err := service.Snapshot(); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	store, err = storage.Open(store.Directory())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	recovered, err := Open(store)
	if err != nil {
		t.Fatal(err)
	}
	item, err := recovered.Get("m-1")
	if err != nil || item.Status != "closed" {
		t.Fatalf("recovery failed: %+v %v", item, err)
	}
}

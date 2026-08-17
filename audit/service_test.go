package audit

import (
	"testing"
	"time"

	"transitmanifest/domain"
	"transitmanifest/storage"
)

func TestAuditListsVerifiedAttributes(t *testing.T) {
	store, _ := storage.Open(t.TempDir())
	defer store.Close()
	event, _ := domain.NewEvent(domain.EventManifestClosed, "m-1", "operator", time.Now(), domain.ManifestClosedData{Reason: "ready"})
	store.Append(event)
	entries, err := New(store).List("m-1", time.Time{}, time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Attributes["reason"] != "ready" {
		t.Fatalf("audit attributes missing: %+v", entries)
	}
	if err := New(store).Verify(); err != nil {
		t.Fatal(err)
	}
}

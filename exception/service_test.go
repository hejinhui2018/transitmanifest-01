package exception

import (
	"testing"
	"time"

	"transitmanifest/domain"
	"transitmanifest/manifest"
	"transitmanifest/storage"
)

func TestRecordAndResolveException(t *testing.T) {
	store, _ := storage.Open(t.TempDir())
	defer store.Close()
	manifests, _ := manifest.Open(store)
	manifests.Create("m-1", "t-1", "plate-1", "a", "b", "operator", time.Now())
	service := New(manifests)
	item, err := service.Record("m-1", "pkg-1", domain.ExceptionDamaged, "wet carton", "a", "worker", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if len(service.List(true)) != 1 {
		t.Fatal("open exception missing")
	}
	if err := service.Resolve(item.ID, "supervisor", time.Now()); err != nil {
		t.Fatal(err)
	}
	if len(service.List(true)) != 0 {
		t.Fatal("resolved exception remains open")
	}
}

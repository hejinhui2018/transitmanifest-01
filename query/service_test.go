package query

import (
	"testing"
	"time"

	"transitmanifest/manifest"
	"transitmanifest/storage"
)

func TestManifestAndPackageQueries(t *testing.T) {
	store, _ := storage.Open(t.TempDir())
	defer store.Close()
	manifests, _ := manifest.Open(store)
	manifests.Create("m-1", "t-1", "plate-1", "a", "b", "operator", time.Now())
	service := New(manifests)
	items := service.Manifests(ManifestFilter{Status: "open", Origin: "a"})
	if len(items) != 1 || items[0].ID != "m-1" {
		t.Fatalf("manifest query failed: %+v", items)
	}
	if history := service.PackageHistory("m-1", "pkg-1"); len(history) != 0 {
		t.Fatalf("unexpected package history: %+v", history)
	}
}

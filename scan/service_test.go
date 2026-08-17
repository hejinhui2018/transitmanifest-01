package scan

import (
	"errors"
	"testing"
	"time"

	"transitmanifest/domain"
	"transitmanifest/manifest"
	"transitmanifest/storage"
)

func TestSubmitIsIdempotentAndReportsPackageDuplicate(t *testing.T) {
	store, _ := storage.Open(t.TempDir())
	defer store.Close()
	manifests, _ := manifest.Open(store)
	manifests.Create("m-1", "t-1", "plate-1", "a", "b", "operator", time.Now())
	service := New(manifests)
	input := domain.PackageScan{ScanID: "scan-1", ManifestID: "m-1", PackageID: "pkg-1", Station: "a", Operator: "operator", Operation: domain.ScanUnload, ScannedAt: time.Now()}
	if _, exception, err := service.Submit(input); err != nil || exception != nil {
		t.Fatalf("first scan failed: %v %v", err, exception)
	}
	if _, exception, err := service.Submit(input); !errors.Is(err, domain.ErrDuplicateScan) || exception != nil {
		t.Fatalf("duplicate was accepted: %v %v", err, exception)
	}
	second := input
	second.ScanID = "scan-2"
	_, exception, err := service.Submit(second)
	if err != nil || exception == nil {
		t.Fatalf("package duplicate not classified: %v %v", err, exception)
	}
	item, _ := manifests.Get("m-1")
	if item.UnloadCount != 2 || item.ExceptionCount != 1 {
		t.Fatalf("counts not updated: %+v", item)
	}
}

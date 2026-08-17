package report

import (
	"testing"
	"time"

	"transitmanifest/domain"
	"transitmanifest/manifest"
	"transitmanifest/scan"
	"transitmanifest/storage"
)

func TestDailyCountsOperationsAndExceptions(t *testing.T) {
	day := time.Date(2026, 8, 18, 9, 0, 0, 0, time.UTC)
	store, _ := storage.Open(t.TempDir())
	defer store.Close()
	manifests, _ := manifest.Open(store)
	manifests.Create("m-1", "t-1", "plate-1", "a", "b", "operator", day)
	service := scan.New(manifests)
	base := domain.PackageScan{ManifestID: "m-1", PackageID: "pkg-1", Station: "a", Operator: "operator", Operation: domain.ScanUnload, ScannedAt: day}
	base.ScanID = "scan-1"
	service.Submit(base)
	base.ScanID = "scan-2"
	service.Submit(base)
	result, err := New(store).Daily(day)
	if err != nil {
		t.Fatal(err)
	}
	if result.Unloaded != 2 || result.Exceptions != 1 || result.ByExceptionKind[string(domain.ExceptionDuplicatePackage)] != 1 {
		t.Fatalf("unexpected report: %+v", result)
	}
}

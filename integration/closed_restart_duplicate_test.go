package integration

import (
	"errors"
	"testing"
	"time"

	"transitmanifest/app"
	"transitmanifest/domain"
)

func TestClosedManifestDuplicateScanSurvivesRestart(t *testing.T) {
	directory := t.TempDir()
	day := time.Date(2026, 8, 18, 10, 0, 0, 0, time.UTC)
	config := app.Config{DataDir: directory, Actor: "station-a"}
	first, err := app.Open(config)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := first.Manifests.Create("m-closed", "trip-7", "plate-7", "hub-a", "hub-b", "dispatcher", day); err != nil {
		t.Fatal(err)
	}
	scan := domain.PackageScan{ScanID: "scan-unload-1", ManifestID: "m-closed", PackageID: "pkg-77", Station: "hub-b", Operator: "dock-7", Operation: domain.ScanUnload, ScannedAt: day.Add(time.Hour)}
	if _, exception, err := first.Scans.Submit(scan); err != nil || exception != nil {
		t.Fatalf("first unload failed: %v %v", err, exception)
	}
	if _, err := first.Manifests.Close("m-closed", "departed", "dispatcher", day.Add(2*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}

	second, err := app.Open(config)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	before, err := second.Reports.Daily(day)
	if err != nil {
		t.Fatal(err)
	}
	if before.Unloaded != 1 || before.Exceptions != 0 {
		t.Fatalf("unexpected precondition report: %+v", before)
	}
	if _, exception, err := second.Scans.Submit(scan); !errors.Is(err, domain.ErrDuplicateScan) || exception != nil {
		t.Fatalf("replayed scan was not rejected: %v %v", err, exception)
	}
	after, err := second.Reports.Daily(day)
	if err != nil {
		t.Fatal(err)
	}
	if after.Unloaded != 1 || after.Exceptions != 0 {
		t.Fatalf("duplicate polluted report: %+v", after)
	}
}

package domain

import (
	"testing"
	"time"
)

func TestValidateManifestAndScan(t *testing.T) {
	manifest := Manifest{ID: "m-1", TripID: "trip-1", VehiclePlate: "plate-1", Origin: "hub-a", Destination: "hub-b"}
	if err := ValidateManifest(manifest); err != nil {
		t.Fatalf("valid manifest rejected: %v", err)
	}
	scan := PackageScan{ScanID: "scan-1", ManifestID: "m-1", PackageID: "pkg-1", Station: "hub-a", Operator: "worker-1", Operation: ScanLoad, ScannedAt: time.Now()}
	if err := ValidateScan(scan); err != nil {
		t.Fatalf("valid scan rejected: %v", err)
	}
	scan.Operation = "move"
	if !IsValidation(ValidateScan(scan)) {
		t.Fatal("invalid operation did not return validation error")
	}
}

package scan

import (
	"strings"
	"time"

	"transitmanifest/domain"
)

func Normalize(scan domain.PackageScan) domain.PackageScan {
	scan.ScanID = strings.TrimSpace(scan.ScanID)
	scan.ManifestID = strings.TrimSpace(scan.ManifestID)
	scan.PackageID = strings.TrimSpace(scan.PackageID)
	scan.Station = strings.TrimSpace(scan.Station)
	scan.Operator = strings.TrimSpace(scan.Operator)
	scan.ScannedAt = scan.ScannedAt.UTC()
	return scan
}

func Prepare(scan domain.PackageScan, now time.Time) (domain.PackageScan, error) {
	scan = Normalize(scan)
	if scan.ScannedAt.IsZero() {
		scan.ScannedAt = now.UTC()
	}
	if err := domain.ValidateScan(scan); err != nil {
		return domain.PackageScan{}, err
	}
	return scan, nil
}

func IsUnload(scan domain.PackageScan) bool { return scan.Operation == domain.ScanUnload }

func IsLoad(scan domain.PackageScan) bool { return scan.Operation == domain.ScanLoad }

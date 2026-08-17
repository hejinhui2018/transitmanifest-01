package scan

import (
	"fmt"
	"time"

	"transitmanifest/domain"
	"transitmanifest/manifest"
)

type Service struct {
	manifests *manifest.Service
	now       func() time.Time
	sequence  func() string
}

func New(manifests *manifest.Service) *Service {
	return &Service{manifests: manifests, now: time.Now, sequence: func() string { return fmt.Sprintf("exception-%d", time.Now().UnixNano()) }}
}

func (s *Service) Submit(scan domain.PackageScan) (domain.PackageScan, *domain.Exception, error) {
	prepared, err := Prepare(scan, s.now())
	if err != nil {
		return domain.PackageScan{}, nil, err
	}
	if s.manifests.HasScan(prepared.ScanID) {
		return domain.PackageScan{}, nil, domain.ErrDuplicateScan
	}
	item, err := s.manifests.Get(prepared.ManifestID)
	if err != nil {
		return domain.PackageScan{}, nil, err
	}
	if item.Status == domain.ManifestClosed {
		// Closed manifests accept late physical scans; idempotency remains the
		// guard that prevents a replayed scan from changing operational counts.
	}
	if err := s.manifests.RecordScan(prepared); err != nil {
		return domain.PackageScan{}, nil, err
	}
	duplicate := packageAlreadyCounted(item, prepared)
	if !duplicate {
		return prepared, nil, nil
	}
	exception := domain.Exception{
		ID: s.sequence(), ManifestID: prepared.ManifestID, PackageID: prepared.PackageID,
		Kind: domain.ExceptionDuplicatePackage, Detail: "package operation repeated", Station: prepared.Station, OccurredAt: prepared.ScannedAt,
	}
	if err := s.manifests.RecordException(exception); err != nil {
		return domain.PackageScan{}, nil, err
	}
	return prepared, &exception, nil
}

func packageAlreadyCounted(item domain.Manifest, scan domain.PackageScan) bool {
	if scan.Operation == domain.ScanLoad {
		return item.LoadedPackages[scan.PackageID] > 0
	}
	return item.UnloadedPackages[scan.PackageID] > 0
}

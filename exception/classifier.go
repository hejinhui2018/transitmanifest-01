package exception

import (
	"fmt"
	"strings"

	"transitmanifest/domain"
)

type Classifier struct {
	RoutePolicy func(routeID, station string) bool
}

func (c Classifier) ForScan(item domain.Manifest, scan domain.PackageScan) (domain.ExceptionKind, bool) {
	if c.RoutePolicy != nil && item.RouteID != "" && !c.RoutePolicy(item.RouteID, scan.Station) {
		return domain.ExceptionWrongRoute, true
	}
	if scan.Operation == domain.ScanLoad && item.LoadedPackages[scan.PackageID] > 0 {
		return domain.ExceptionDuplicatePackage, true
	}
	if scan.Operation == domain.ScanUnload && item.UnloadedPackages[scan.PackageID] > 0 {
		return domain.ExceptionDuplicatePackage, true
	}
	return "", false
}

func ParseKind(value string) (domain.ExceptionKind, error) {
	kind := domain.ExceptionKind(strings.ToLower(strings.TrimSpace(value)))
	switch kind {
	case domain.ExceptionDuplicatePackage, domain.ExceptionWrongRoute, domain.ExceptionDamaged, domain.ExceptionMissing:
		return kind, nil
	default:
		return "", fmt.Errorf("unknown exception kind %q", value)
	}
}

func IsOpen(item domain.Exception) bool { return item.ResolvedAt == nil }

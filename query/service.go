package query

import (
	"sort"

	"transitmanifest/domain"
	"transitmanifest/manifest"
)

type Service struct{ manifests *manifest.Service }

func New(manifests *manifest.Service) *Service { return &Service{manifests: manifests} }

func (s *Service) Manifests(filter ManifestFilter) []domain.Manifest {
	items := s.manifests.List()
	filtered := items[:0]
	for _, item := range items {
		if Matches(item, filter) {
			filtered = append(filtered, item)
		}
	}
	SortByCreated(filtered)
	return filtered
}

func (s *Service) PackageHistory(manifestID, packageID string) []domain.PackageScan {
	items := s.manifests.Scans()
	result := items[:0]
	for _, item := range items {
		if item.ManifestID == manifestID && (packageID == "" || item.PackageID == packageID) {
			result = append(result, item)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ScannedAt.Before(result[j].ScannedAt) })
	return result
}

func (s *Service) ExceptionHistory(manifestID string) []domain.Exception {
	items := s.manifests.Exceptions()
	result := items[:0]
	for _, item := range items {
		if manifestID == "" || item.ManifestID == manifestID {
			result = append(result, item)
		}
	}
	return result
}

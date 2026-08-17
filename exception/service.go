package exception

import (
	"sort"
	"time"

	"transitmanifest/domain"
	"transitmanifest/manifest"
)

type Service struct{ manifests *manifest.Service }

func New(manifests *manifest.Service) *Service { return &Service{manifests: manifests} }

func (s *Service) List(openOnly bool) []domain.Exception {
	items := s.manifests.Exceptions()
	if openOnly {
		filtered := items[:0]
		for _, item := range items {
			if IsOpen(item) {
				filtered = append(filtered, item)
			}
		}
		items = filtered
	}
	sort.Slice(items, func(i, j int) bool { return items[i].OccurredAt.Before(items[j].OccurredAt) })
	return items
}

func (s *Service) Record(manifestID, packageID string, kind domain.ExceptionKind, detail, station, actor string, at time.Time) (domain.Exception, error) {
	if at.IsZero() {
		at = time.Now()
	}
	exception := domain.Exception{ID: actor + "-" + at.Format("20060102150405.000000000"), ManifestID: manifestID, PackageID: packageID, Kind: kind, Detail: detail, Station: station, OccurredAt: at}
	if err := s.manifests.RecordException(exception); err != nil {
		return domain.Exception{}, err
	}
	return exception, nil
}

func (s *Service) Resolve(id, actor string, at time.Time) error {
	return s.manifests.ResolveException(id, actor, at)
}

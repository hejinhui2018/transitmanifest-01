package report

import (
	"encoding/json"
	"time"

	"transitmanifest/domain"
	"transitmanifest/storage"
)

type Service struct{ store *storage.Store }

func New(store *storage.Store) *Service { return &Service{store: store} }

func (s *Service) Daily(day time.Time) (domain.DailyReport, error) {
	result := domain.DailyReport{Date: day.Format("2006-01-02"), ByExceptionKind: map[string]int{}, ByStation: map[string]int{}}
	records, err := s.store.Records()
	if err != nil {
		return result, err
	}
	for _, record := range records {
		if record.OccurredAt.Format("2006-01-02") != result.Date {
			continue
		}
		switch record.Type {
		case domain.EventManifestCreated:
			result.CreatedManifests++
		case domain.EventManifestClosed:
			result.ClosedManifests++
		case domain.EventPackageScanned:
			var scan domain.PackageScan
			if err := decode(record, &scan); err != nil {
				return result, err
			}
			if scan.Operation == domain.ScanLoad {
				result.Loaded++
			} else {
				result.Unloaded++
			}
			result.ByStation[scan.Station]++
		case domain.EventHandoffSigned:
			result.Handoffs++
		case domain.EventExceptionRecorded:
			var exception domain.Exception
			if err := decode(record, &exception); err != nil {
				return result, err
			}
			result.Exceptions++
			result.ByExceptionKind[string(exception.Kind)]++
		}
	}
	return result, nil
}

func decode(record storage.Record, target any) error {
	return json.Unmarshal(record.Data, target)
}

package scan

import (
	"sort"

	"transitmanifest/domain"
)

type BatchResult struct {
	Accepted   []domain.PackageScan
	Exceptions []domain.Exception
	Rejected   []error
}

func (s *Service) SubmitBatch(scans []domain.PackageScan) BatchResult {
	result := BatchResult{Accepted: make([]domain.PackageScan, 0, len(scans)), Exceptions: make([]domain.Exception, 0), Rejected: make([]error, 0)}
	for _, input := range scans {
		accepted, exception, err := s.Submit(input)
		if err != nil {
			result.Rejected = append(result.Rejected, err)
			continue
		}
		result.Accepted = append(result.Accepted, accepted)
		if exception != nil {
			result.Exceptions = append(result.Exceptions, *exception)
		}
	}
	sort.Slice(result.Accepted, func(i, j int) bool { return result.Accepted[i].ScannedAt.Before(result.Accepted[j].ScannedAt) })
	return result
}

func CountOperations(scans []domain.PackageScan) (loads, unloads int) {
	for _, item := range scans {
		if item.Operation == domain.ScanLoad {
			loads++
		}
		if item.Operation == domain.ScanUnload {
			unloads++
		}
	}
	return loads, unloads
}

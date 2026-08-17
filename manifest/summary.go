package manifest

import (
	"fmt"
	"sort"

	"transitmanifest/domain"
)

type Summary struct {
	ManifestID       string `json:"manifest_id"`
	Status           string `json:"status"`
	LoadCount        int    `json:"load_count"`
	UnloadCount      int    `json:"unload_count"`
	LoadedDistinct   int    `json:"loaded_distinct"`
	UnloadedDistinct int    `json:"unloaded_distinct"`
	ExceptionCount   int    `json:"exception_count"`
	HasHandoff       bool   `json:"has_handoff"`
}

func Summarize(item domain.Manifest) Summary {
	return Summary{
		ManifestID: item.ID, Status: string(item.Status), LoadCount: item.LoadCount,
		UnloadCount: item.UnloadCount, LoadedDistinct: len(item.LoadedPackages),
		UnloadedDistinct: len(item.UnloadedPackages), ExceptionCount: item.ExceptionCount,
		HasHandoff: item.Handoff != nil,
	}
}

func (s *Service) Summaries() []Summary {
	items := s.List()
	result := make([]Summary, 0, len(items))
	for _, item := range items {
		result = append(result, Summarize(item))
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ManifestID < result[j].ManifestID })
	return result
}

// ValidateInvariants is used by operations diagnostics and by callers after
// recovery. It checks aggregate counters against the durable package maps.
func (s *Service) ValidateInvariants() error {
	for _, item := range s.List() {
		loaded, unloaded := 0, 0
		for _, count := range item.LoadedPackages {
			loaded += count
		}
		for _, count := range item.UnloadedPackages {
			unloaded += count
		}
		if loaded != item.LoadCount || unloaded != item.UnloadCount {
			return fmt.Errorf("manifest %s counters disagree with package maps", item.ID)
		}
		if item.Status == domain.ManifestClosed && item.ClosedAt == nil {
			return fmt.Errorf("manifest %s is closed without close time", item.ID)
		}
	}
	return nil
}

func PackageBalance(item domain.Manifest, packageID string) int {
	return item.LoadedPackages[packageID] - item.UnloadedPackages[packageID]
}

func OpenItems(items []domain.Manifest) []domain.Manifest {
	result := make([]domain.Manifest, 0)
	for _, item := range items {
		if item.Status == domain.ManifestOpen {
			result = append(result, item.Clone())
		}
	}
	return result
}

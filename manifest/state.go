package manifest

import (
	"sort"

	"transitmanifest/domain"
)

type persistedState struct {
	Manifests  map[string]domain.Manifest    `json:"manifests"`
	Scans      map[string]domain.PackageScan `json:"scans"`
	Routes     map[string]domain.Route       `json:"routes"`
	Exceptions map[string]domain.Exception   `json:"exceptions"`
}

func newState() persistedState {
	return persistedState{
		Manifests:  make(map[string]domain.Manifest),
		Scans:      make(map[string]domain.PackageScan),
		Routes:     make(map[string]domain.Route),
		Exceptions: make(map[string]domain.Exception),
	}
}

func (s persistedState) clone() persistedState {
	copyState := newState()
	for id, item := range s.Manifests {
		copyState.Manifests[id] = item.Clone()
	}
	for id, item := range s.Scans {
		copyState.Scans[id] = item
	}
	for id, item := range s.Routes {
		item.Via = append([]string(nil), item.Via...)
		copyState.Routes[id] = item
	}
	for id, item := range s.Exceptions {
		copyState.Exceptions[id] = item
	}
	return copyState
}

func sortedManifests(source map[string]domain.Manifest) []domain.Manifest {
	result := make([]domain.Manifest, 0, len(source))
	for _, item := range source {
		result = append(result, item.Clone())
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

func sortedExceptions(source map[string]domain.Exception) []domain.Exception {
	result := make([]domain.Exception, 0, len(source))
	for _, item := range source {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].OccurredAt.Equal(result[j].OccurredAt) {
			return result[i].ID < result[j].ID
		}
		return result[i].OccurredAt.Before(result[j].OccurredAt)
	})
	return result
}

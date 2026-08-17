package query

import (
	"sort"
	"strings"

	"transitmanifest/domain"
)

type ManifestFilter struct {
	Status      domain.ManifestStatus
	Origin      string
	Destination string
	Vehicle     string
	RouteID     string
}

func Matches(item domain.Manifest, filter ManifestFilter) bool {
	if filter.Status != "" && item.Status != filter.Status {
		return false
	}
	if filter.Origin != "" && !strings.EqualFold(item.Origin, filter.Origin) {
		return false
	}
	if filter.Destination != "" && !strings.EqualFold(item.Destination, filter.Destination) {
		return false
	}
	if filter.Vehicle != "" && !strings.EqualFold(item.VehiclePlate, filter.Vehicle) {
		return false
	}
	if filter.RouteID != "" && item.RouteID != filter.RouteID {
		return false
	}
	return true
}

func SortByCreated(items []domain.Manifest) {
	sort.SliceStable(items, func(i, j int) bool { return items[i].CreatedAt.Before(items[j].CreatedAt) })
}

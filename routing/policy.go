package routing

import (
	"fmt"
	"strings"

	"transitmanifest/domain"
)

type Policy struct {
	Routes map[string]domain.Route
}

func NewPolicy(routes []domain.Route) Policy {
	result := Policy{Routes: make(map[string]domain.Route, len(routes))}
	for _, route := range routes {
		result.Routes[route.ID] = route
	}
	return result
}

func (p Policy) Resolve(origin, destination string) (domain.Route, error) {
	for _, route := range p.Routes {
		if strings.EqualFold(route.Origin, origin) && strings.EqualFold(route.Destination, destination) && route.Enabled {
			return route, nil
		}
	}
	return domain.Route{}, fmt.Errorf("%w: %s -> %s", domain.ErrRouteNotFound, origin, destination)
}

func (p Policy) Allows(routeID, station string) bool {
	route, ok := p.Routes[routeID]
	if !ok || !route.Enabled {
		return false
	}
	if station == route.Origin || station == route.Destination {
		return true
	}
	for _, stop := range route.Via {
		if stop == station {
			return true
		}
	}
	return false
}

package routing

import (
	"sync"
	"time"

	"transitmanifest/domain"
	"transitmanifest/manifest"
)

type Service struct {
	mu        sync.RWMutex
	manifests *manifest.Service
	routes    map[string]domain.Route
}

func New(manifests *manifest.Service) *Service {
	routes := make(map[string]domain.Route)
	for _, route := range manifests.Routes() {
		routes[route.ID] = route
	}
	return &Service{manifests: manifests, routes: routes}
}

func (s *Service) Register(route domain.Route) error {
	if err := domain.ValidateRoute(route); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.routes[route.ID] = route
	return nil
}

func (s *Service) Assign(manifestID, routeID, actor string) error {
	s.mu.RLock()
	route, ok := s.routes[routeID]
	s.mu.RUnlock()
	if !ok {
		return domain.ErrRouteNotFound
	}
	if !route.Enabled {
		return &domain.ConflictError{ManifestID: manifestID, Operation: "assign route", Reason: "route disabled"}
	}
	return s.manifests.AssignRoute(manifestID, route, actor, time.Now())
}

func (s *Service) Policy() Policy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	routes := make([]domain.Route, 0, len(s.routes))
	for _, route := range s.routes {
		routes = append(routes, route)
	}
	return NewPolicy(routes)
}

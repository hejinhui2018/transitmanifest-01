package manifest

import (
	"fmt"
	"time"

	"transitmanifest/domain"
)

func (s *Service) Create(id, tripID, plate, origin, destination, actor string, at time.Time) (domain.Manifest, error) {
	if at.IsZero() {
		at = s.now()
	}
	item := domain.Manifest{ID: id, TripID: tripID, VehiclePlate: plate, Origin: origin, Destination: destination, Status: domain.ManifestOpen}
	if err := domain.ValidateManifest(item); err != nil {
		return domain.Manifest{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.state.Manifests[id]; exists {
		return domain.Manifest{}, domain.ErrManifestExists
	}
	event, err := domain.NewEvent(domain.EventManifestCreated, id, actor, at, domain.ManifestCreated{TripID: tripID, VehiclePlate: plate, Origin: origin, Destination: destination})
	if err != nil {
		return domain.Manifest{}, err
	}
	if _, err := s.store.Append(event); err != nil {
		return domain.Manifest{}, err
	}
	if err := applyEvent(&s.state, event); err != nil {
		return domain.Manifest{}, err
	}
	return s.state.Manifests[id].Clone(), nil
}

func (s *Service) AssignRoute(manifestID string, route domain.Route, actor string, at time.Time) error {
	if at.IsZero() {
		at = s.now()
	}
	if err := domain.ValidateRoute(route); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.state.Manifests[manifestID]
	if !ok {
		return domain.ErrManifestNotFound
	}
	if item.Status == domain.ManifestClosed {
		return &domain.ConflictError{ManifestID: manifestID, Operation: "assign a route", Reason: domain.ErrManifestClosed.Error()}
	}
	s.state.Routes[route.ID] = route
	event, err := domain.NewEvent(domain.EventRouteAssigned, manifestID, actor, at, domain.RouteAssigned{RouteID: route.ID})
	if err != nil {
		return err
	}
	if _, err := s.store.Append(event); err != nil {
		return err
	}
	return applyEvent(&s.state, event)
}

func (s *Service) Close(id, reason, actor string, at time.Time) (domain.Manifest, error) {
	if at.IsZero() {
		at = s.now()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.state.Manifests[id]
	if !ok {
		return domain.Manifest{}, domain.ErrManifestNotFound
	}
	if item.Status == domain.ManifestClosed {
		return domain.Manifest{}, &domain.ConflictError{ManifestID: id, Operation: "close", Reason: "already closed"}
	}
	event, err := domain.NewEvent(domain.EventManifestClosed, id, actor, at, domain.ManifestClosedData{Reason: reason})
	if err != nil {
		return domain.Manifest{}, err
	}
	if _, err := s.store.Append(event); err != nil {
		return domain.Manifest{}, err
	}
	if err := applyEvent(&s.state, event); err != nil {
		return domain.Manifest{}, err
	}
	return s.state.Manifests[id].Clone(), nil
}

func (s *Service) RecordScan(scan domain.PackageScan) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.state.Manifests[scan.ManifestID]; !ok {
		return domain.ErrManifestNotFound
	}
	if _, ok := s.state.Scans[scan.ScanID]; ok {
		return domain.ErrDuplicateScan
	}
	event, err := domain.NewEvent(domain.EventPackageScanned, scan.ManifestID, scan.Operator, scan.ScannedAt, scan)
	if err != nil {
		return err
	}
	if _, err := s.store.Append(event); err != nil {
		return err
	}
	return applyEvent(&s.state, event)
}

func (s *Service) RecordException(exception domain.Exception) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.state.Manifests[exception.ManifestID]; !ok {
		return domain.ErrManifestNotFound
	}
	if exception.ID == "" {
		return &domain.FieldError{Field: "exception_id", Problem: "is required"}
	}
	event, err := domain.NewEvent(domain.EventExceptionRecorded, exception.ManifestID, exception.Station, exception.OccurredAt, exception)
	if err != nil {
		return err
	}
	if _, err := s.store.Append(event); err != nil {
		return err
	}
	return applyEvent(&s.state, event)
}

func (s *Service) SignHandoff(receipt domain.HandoffReceipt) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.state.Manifests[receipt.ManifestID]
	if !ok {
		return domain.ErrManifestNotFound
	}
	if item.Handoff != nil {
		return domain.ErrHandoffExists
	}
	if receipt.SignedAt.IsZero() {
		receipt.SignedAt = s.now()
	}
	event, err := domain.NewEvent(domain.EventHandoffSigned, receipt.ManifestID, receipt.Signer, receipt.SignedAt, receipt)
	if err != nil {
		return err
	}
	if _, err := s.store.Append(event); err != nil {
		return err
	}
	return applyEvent(&s.state, event)
}

func (s *Service) ResolveException(exceptionID, actor string, at time.Time) error {
	if at.IsZero() {
		at = s.now()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.state.Exceptions[exceptionID]
	if !ok {
		return fmt.Errorf("exception %s not found", exceptionID)
	}
	if item.ResolvedAt != nil {
		return fmt.Errorf("exception %s already resolved", exceptionID)
	}
	event, err := domain.NewEvent(domain.EventExceptionResolved, item.ManifestID, actor, at, domain.ExceptionResolvedData{ExceptionID: exceptionID})
	if err != nil {
		return err
	}
	if _, err := s.store.Append(event); err != nil {
		return err
	}
	return applyEvent(&s.state, event)
}

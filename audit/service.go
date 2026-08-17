package audit

import (
	"encoding/json"
	"fmt"
	"time"

	"transitmanifest/domain"
	"transitmanifest/storage"
)

type Service struct{ store *storage.Store }

func New(store *storage.Store) *Service { return &Service{store: store} }

func (s *Service) List(aggregateID string, from, to time.Time) ([]domain.AuditEntry, error) {
	records, err := s.store.Records()
	if err != nil {
		return nil, err
	}
	result := make([]domain.AuditEntry, 0)
	for _, record := range records {
		if aggregateID != "" && record.AggregateID != aggregateID {
			continue
		}
		if !from.IsZero() && record.OccurredAt.Before(from) {
			continue
		}
		if !to.IsZero() && record.OccurredAt.After(to) {
			continue
		}
		attributes := map[string]any{}
		_ = json.Unmarshal(record.Data, &attributes)
		result = append(result, domain.AuditEntry{Sequence: record.Sequence, EventType: record.Type, Aggregate: record.AggregateID, Actor: record.Actor, OccurredAt: record.OccurredAt, Attributes: attributes})
	}
	return result, nil
}

func (s *Service) Verify() error {
	_, err := s.store.Records()
	if err != nil {
		return fmt.Errorf("audit verification failed: %w", err)
	}
	return nil
}

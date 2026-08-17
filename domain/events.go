package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	EventManifestCreated   = "manifest.created"
	EventRouteAssigned     = "manifest.route_assigned"
	EventPackageScanned    = "package.scanned"
	EventManifestClosed    = "manifest.closed"
	EventHandoffSigned     = "handoff.signed"
	EventExceptionRecorded = "exception.recorded"
	EventExceptionResolved = "exception.resolved"
)

type Event struct {
	Type        string          `json:"type"`
	AggregateID string          `json:"aggregate_id"`
	Actor       string          `json:"actor"`
	OccurredAt  time.Time       `json:"occurred_at"`
	Data        json.RawMessage `json:"data"`
}

func NewEvent(kind, aggregateID, actor string, occurredAt time.Time, payload any) (Event, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return Event{}, fmt.Errorf("marshal %s event: %w", kind, err)
	}
	return Event{
		Type:        kind,
		AggregateID: aggregateID,
		Actor:       actor,
		OccurredAt:  occurredAt.UTC(),
		Data:        data,
	}, nil
}

func DecodeEvent[T any](event Event) (T, error) {
	var value T
	if err := json.Unmarshal(event.Data, &value); err != nil {
		return value, fmt.Errorf("decode %s event: %w", event.Type, err)
	}
	return value, nil
}

type ManifestCreated struct {
	TripID       string `json:"trip_id"`
	VehiclePlate string `json:"vehicle_plate"`
	Origin       string `json:"origin"`
	Destination  string `json:"destination"`
}

type RouteAssigned struct {
	RouteID string `json:"route_id"`
}

type ManifestClosedData struct {
	Reason string `json:"reason,omitempty"`
}

type ExceptionResolvedData struct {
	ExceptionID string `json:"exception_id"`
}

// ValidateType prevents unrecognized records from silently entering the
// audit trail. Unknown historical event types are treated as corruption.
func ValidateType(kind string) error {
	switch kind {
	case EventManifestCreated, EventRouteAssigned, EventPackageScanned,
		EventManifestClosed, EventHandoffSigned, EventExceptionRecorded,
		EventExceptionResolved:
		return nil
	default:
		return fmt.Errorf("%w: unknown event type %q", ErrCorruptLog, kind)
	}
}

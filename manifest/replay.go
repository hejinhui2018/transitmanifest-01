package manifest

import (
	"fmt"

	"transitmanifest/domain"
	"transitmanifest/storage"
)

func recoverState(store *storage.Store) (persistedState, error) {
	state := newState()
	snapshot, records, err := store.Recover(&state)
	if err != nil {
		return state, err
	}
	// The snapshot payload already restores the full persisted state,
	// including the scan-id idempotency index. Only records appended after
	// the snapshot sequence are replayed below, so the index must be kept to
	// reject duplicate scans that were already applied before the restart.
	_ = snapshot
	for _, record := range records {
		if err := applyRecord(&state, record); err != nil {
			return state, fmt.Errorf("replay sequence %d: %w", record.Sequence, err)
		}
	}
	return state, nil
}

func applyRecord(state *persistedState, record storage.Record) error {
	return applyEvent(state, record.Event())
}

func applyEvent(state *persistedState, event domain.Event) error {
	switch event.Type {
	case domain.EventManifestCreated:
		data, err := domain.DecodeEvent[domain.ManifestCreated](event)
		if err != nil {
			return err
		}
		state.Manifests[event.AggregateID] = domain.Manifest{
			ID: event.AggregateID, TripID: data.TripID, VehiclePlate: data.VehiclePlate,
			Origin: data.Origin, Destination: data.Destination, Status: domain.ManifestOpen,
			CreatedAt: event.OccurredAt, LoadedPackages: map[string]int{}, UnloadedPackages: map[string]int{},
		}
	case domain.EventRouteAssigned:
		data, err := domain.DecodeEvent[domain.RouteAssigned](event)
		if err != nil {
			return err
		}
		item, ok := state.Manifests[event.AggregateID]
		if !ok {
			return fmt.Errorf("%w: %s", domain.ErrManifestNotFound, event.AggregateID)
		}
		item.RouteID = data.RouteID
		state.Manifests[event.AggregateID] = item
	case domain.EventPackageScanned:
		scan, err := domain.DecodeEvent[domain.PackageScan](event)
		if err != nil {
			return err
		}
		item, ok := state.Manifests[scan.ManifestID]
		if !ok {
			return fmt.Errorf("%w: %s", domain.ErrManifestNotFound, scan.ManifestID)
		}
		if item.LoadedPackages == nil {
			item.LoadedPackages = map[string]int{}
		}
		if item.UnloadedPackages == nil {
			item.UnloadedPackages = map[string]int{}
		}
		if scan.Operation == domain.ScanLoad {
			item.LoadCount++
			item.LoadedPackages[scan.PackageID]++
		}
		if scan.Operation == domain.ScanUnload {
			item.UnloadCount++
			item.UnloadedPackages[scan.PackageID]++
		}
		state.Manifests[scan.ManifestID] = item
		state.Scans[scan.ScanID] = scan
	case domain.EventManifestClosed:
		item, ok := state.Manifests[event.AggregateID]
		if !ok {
			return fmt.Errorf("%w: %s", domain.ErrManifestNotFound, event.AggregateID)
		}
		item.Status = domain.ManifestClosed
		closed := event.OccurredAt
		item.ClosedAt = &closed
		state.Manifests[event.AggregateID] = item
	case domain.EventHandoffSigned:
		receipt, err := domain.DecodeEvent[domain.HandoffReceipt](event)
		if err != nil {
			return err
		}
		item, ok := state.Manifests[event.AggregateID]
		if !ok {
			return fmt.Errorf("%w: %s", domain.ErrManifestNotFound, event.AggregateID)
		}
		item.Handoff = &receipt
		state.Manifests[event.AggregateID] = item
	case domain.EventExceptionRecorded:
		exception, err := domain.DecodeEvent[domain.Exception](event)
		if err != nil {
			return err
		}
		state.Exceptions[exception.ID] = exception
		item := state.Manifests[exception.ManifestID]
		item.ExceptionCount++
		state.Manifests[exception.ManifestID] = item
	case domain.EventExceptionResolved:
		data, err := domain.DecodeEvent[domain.ExceptionResolvedData](event)
		if err != nil {
			return err
		}
		exception, ok := state.Exceptions[data.ExceptionID]
		if !ok {
			return fmt.Errorf("exception %s not found", data.ExceptionID)
		}
		resolved := event.OccurredAt
		exception.ResolvedAt = &resolved
		state.Exceptions[data.ExceptionID] = exception
	default:
		return fmt.Errorf("%w: %s", domain.ErrCorruptLog, event.Type)
	}
	return nil
}

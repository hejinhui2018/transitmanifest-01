package domain

import "time"

type ManifestStatus string

const (
	ManifestOpen   ManifestStatus = "open"
	ManifestClosed ManifestStatus = "closed"
)

type ScanOperation string

const (
	ScanLoad   ScanOperation = "load"
	ScanUnload ScanOperation = "unload"
)

type ExceptionKind string

const (
	ExceptionDuplicatePackage ExceptionKind = "duplicate_package"
	ExceptionWrongRoute       ExceptionKind = "wrong_route"
	ExceptionDamaged          ExceptionKind = "damaged"
	ExceptionMissing          ExceptionKind = "missing"
)

type Manifest struct {
	ID               string          `json:"id"`
	TripID           string          `json:"trip_id"`
	VehiclePlate     string          `json:"vehicle_plate"`
	Origin           string          `json:"origin"`
	Destination      string          `json:"destination"`
	RouteID          string          `json:"route_id,omitempty"`
	Status           ManifestStatus  `json:"status"`
	CreatedAt        time.Time       `json:"created_at"`
	ClosedAt         *time.Time      `json:"closed_at,omitempty"`
	LoadedPackages   map[string]int  `json:"loaded_packages"`
	UnloadedPackages map[string]int  `json:"unloaded_packages"`
	LoadCount        int             `json:"load_count"`
	UnloadCount      int             `json:"unload_count"`
	ExceptionCount   int             `json:"exception_count"`
	Handoff          *HandoffReceipt `json:"handoff,omitempty"`
}

type PackageScan struct {
	ScanID     string        `json:"scan_id"`
	ManifestID string        `json:"manifest_id"`
	PackageID  string        `json:"package_id"`
	Station    string        `json:"station"`
	Operator   string        `json:"operator"`
	Operation  ScanOperation `json:"operation"`
	ScannedAt  time.Time     `json:"scanned_at"`
}

type HandoffReceipt struct {
	ManifestID string    `json:"manifest_id"`
	From       string    `json:"from"`
	To         string    `json:"to"`
	Signer     string    `json:"signer"`
	SignedAt   time.Time `json:"signed_at"`
	Note       string    `json:"note,omitempty"`
}

type Exception struct {
	ID         string        `json:"id"`
	ManifestID string        `json:"manifest_id"`
	PackageID  string        `json:"package_id,omitempty"`
	Kind       ExceptionKind `json:"kind"`
	Detail     string        `json:"detail"`
	Station    string        `json:"station"`
	OccurredAt time.Time     `json:"occurred_at"`
	ResolvedAt *time.Time    `json:"resolved_at,omitempty"`
}

type Route struct {
	ID          string   `json:"id"`
	Origin      string   `json:"origin"`
	Destination string   `json:"destination"`
	Via         []string `json:"via,omitempty"`
	Enabled     bool     `json:"enabled"`
}

type AuditEntry struct {
	Sequence   uint64         `json:"sequence"`
	EventType  string         `json:"event_type"`
	Aggregate  string         `json:"aggregate"`
	Actor      string         `json:"actor"`
	OccurredAt time.Time      `json:"occurred_at"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

type DailyReport struct {
	Date             string         `json:"date"`
	CreatedManifests int            `json:"created_manifests"`
	ClosedManifests  int            `json:"closed_manifests"`
	Loaded           int            `json:"loaded"`
	Unloaded         int            `json:"unloaded"`
	Handoffs         int            `json:"handoffs"`
	Exceptions       int            `json:"exceptions"`
	ByExceptionKind  map[string]int `json:"by_exception_kind"`
	ByStation        map[string]int `json:"by_station"`
}

func (m Manifest) Clone() Manifest {
	m.LoadedPackages = cloneCounts(m.LoadedPackages)
	m.UnloadedPackages = cloneCounts(m.UnloadedPackages)
	if m.ClosedAt != nil {
		v := *m.ClosedAt
		m.ClosedAt = &v
	}
	if m.Handoff != nil {
		v := *m.Handoff
		m.Handoff = &v
	}
	return m
}

func cloneCounts(source map[string]int) map[string]int {
	result := make(map[string]int, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}

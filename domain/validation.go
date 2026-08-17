package domain

import (
	"regexp"
	"strings"
	"time"
)

var identifierPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,63}$`)

func ValidateManifest(m Manifest) error {
	checks := []struct {
		name  string
		value string
	}{
		{"manifest_id", m.ID},
		{"trip_id", m.TripID},
		{"vehicle_plate", m.VehiclePlate},
		{"origin", m.Origin},
		{"destination", m.Destination},
	}
	for _, check := range checks {
		if err := ValidateIdentifier(check.name, check.value); err != nil {
			return err
		}
	}
	if strings.EqualFold(m.Origin, m.Destination) {
		return &FieldError{Field: "destination", Value: m.Destination, Problem: "must differ from origin"}
	}
	return nil
}

func ValidateScan(scan PackageScan) error {
	checks := []struct {
		name  string
		value string
	}{
		{"scan_id", scan.ScanID},
		{"manifest_id", scan.ManifestID},
		{"package_id", scan.PackageID},
		{"station", scan.Station},
		{"operator", scan.Operator},
	}
	for _, check := range checks {
		if err := ValidateIdentifier(check.name, check.value); err != nil {
			return err
		}
	}
	if scan.Operation != ScanLoad && scan.Operation != ScanUnload {
		return &FieldError{Field: "operation", Value: string(scan.Operation), Problem: "must be load or unload"}
	}
	if scan.ScannedAt.IsZero() {
		return &FieldError{Field: "scanned_at", Problem: "is required"}
	}
	if scan.ScannedAt.After(time.Now().Add(24 * time.Hour)) {
		return &FieldError{Field: "scanned_at", Problem: "is too far in the future"}
	}
	return nil
}

func ValidateIdentifier(field, value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return &FieldError{Field: field, Problem: "is required"}
	}
	if !identifierPattern.MatchString(value) {
		return &FieldError{Field: field, Value: value, Problem: "contains unsupported characters or length"}
	}
	return nil
}

func ValidateRoute(route Route) error {
	if err := ValidateIdentifier("route_id", route.ID); err != nil {
		return err
	}
	if err := ValidateIdentifier("origin", route.Origin); err != nil {
		return err
	}
	if err := ValidateIdentifier("destination", route.Destination); err != nil {
		return err
	}
	if strings.EqualFold(route.Origin, route.Destination) {
		return &FieldError{Field: "destination", Value: route.Destination, Problem: "must differ from origin"}
	}
	for _, station := range route.Via {
		if err := ValidateIdentifier("via", station); err != nil {
			return err
		}
	}
	return nil
}

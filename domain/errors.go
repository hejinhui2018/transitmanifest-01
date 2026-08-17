package domain

import (
	"errors"
	"fmt"
)

var (
	ErrManifestNotFound = errors.New("manifest not found")
	ErrManifestExists   = errors.New("manifest already exists")
	ErrManifestClosed   = errors.New("manifest already closed")
	ErrDuplicateScan    = errors.New("duplicate scan")
	ErrInvalidInput     = errors.New("invalid input")
	ErrRouteNotFound    = errors.New("route not found")
	ErrHandoffExists    = errors.New("handoff already signed")
	ErrCorruptLog       = errors.New("corrupt event log")
	ErrCorruptSnapshot  = errors.New("corrupt snapshot")
)

// FieldError preserves enough context for CLI users to correct a command
// without exposing storage or implementation details.
type FieldError struct {
	Field   string
	Value   string
	Problem string
}

func (e *FieldError) Error() string {
	if e.Value == "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Problem)
	}
	return fmt.Sprintf("%s %q: %s", e.Field, e.Value, e.Problem)
}

func (e *FieldError) Unwrap() error {
	return ErrInvalidInput
}

// ConflictError identifies a command that is valid in isolation but cannot
// be applied to the current aggregate state.
type ConflictError struct {
	ManifestID string
	Operation  string
	Reason     string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("manifest %s cannot %s: %s", e.ManifestID, e.Operation, e.Reason)
}

func IsValidation(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

func IsConflict(err error) bool {
	var target *ConflictError
	return errors.As(err, &target)
}

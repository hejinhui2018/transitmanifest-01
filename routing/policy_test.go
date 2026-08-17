package routing

import (
	"testing"

	"transitmanifest/domain"
)

func TestPolicyResolvesEnabledRouteAndStations(t *testing.T) {
	policy := NewPolicy([]domain.Route{{ID: "r-1", Origin: "a", Destination: "c", Via: []string{"b"}, Enabled: true}})
	if _, err := policy.Resolve("a", "c"); err != nil {
		t.Fatal(err)
	}
	if !policy.Allows("r-1", "b") || policy.Allows("r-1", "x") {
		t.Fatal("route station policy incorrect")
	}
}

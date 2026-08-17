package audit

import (
	"sort"

	"transitmanifest/domain"
)

func GroupByType(entries []domain.AuditEntry) map[string]int {
	result := make(map[string]int)
	for _, entry := range entries {
		result[entry.EventType]++
	}
	return result
}

func Timeline(entries []domain.AuditEntry) []domain.AuditEntry {
	result := append([]domain.AuditEntry(nil), entries...)
	sort.SliceStable(result, func(i, j int) bool { return result[i].Sequence < result[j].Sequence })
	return result
}

func Actors(entries []domain.AuditEntry) []string {
	seen := map[string]bool{}
	result := make([]string, 0)
	for _, entry := range entries {
		if !seen[entry.Actor] {
			seen[entry.Actor] = true
			result = append(result, entry.Actor)
		}
	}
	sort.Strings(result)
	return result
}

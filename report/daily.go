package report

import (
	"sort"

	"transitmanifest/domain"
)

func Merge(reports ...domain.DailyReport) domain.DailyReport {
	result := domain.DailyReport{ByExceptionKind: map[string]int{}, ByStation: map[string]int{}}
	for _, report := range reports {
		if result.Date == "" {
			result.Date = report.Date
		}
		result.CreatedManifests += report.CreatedManifests
		result.ClosedManifests += report.ClosedManifests
		result.Loaded += report.Loaded
		result.Unloaded += report.Unloaded
		result.Handoffs += report.Handoffs
		result.Exceptions += report.Exceptions
		for key, count := range report.ByExceptionKind {
			result.ByExceptionKind[key] += count
		}
		for key, count := range report.ByStation {
			result.ByStation[key] += count
		}
	}
	return result
}

func TopStations(report domain.DailyReport, limit int) []string {
	keys := make([]string, 0, len(report.ByStation))
	for key := range report.ByStation {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if report.ByStation[keys[i]] == report.ByStation[keys[j]] {
			return keys[i] < keys[j]
		}
		return report.ByStation[keys[i]] > report.ByStation[keys[j]]
	})
	if limit < len(keys) {
		keys = keys[:limit]
	}
	return keys
}

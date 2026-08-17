package report

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"transitmanifest/domain"
)

func WriteJSON(writer io.Writer, value domain.DailyReport) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		return fmt.Errorf("encode daily report: %w", err)
	}
	return nil
}

func Markdown(report domain.DailyReport) string {
	lines := []string{
		"# TransitManifest Daily Report",
		"",
		"| Date | Created | Closed | Loaded | Unloaded | Handoffs | Exceptions |",
		"|---|---:|---:|---:|---:|---:|---:|",
		fmt.Sprintf("| %s | %d | %d | %d | %d | %d | %d |", report.Date, report.CreatedManifests, report.ClosedManifests, report.Loaded, report.Unloaded, report.Handoffs, report.Exceptions),
		"",
		"## Exceptions",
		"",
	}
	keys := make([]string, 0, len(report.ByExceptionKind))
	for key := range report.ByExceptionKind {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		lines = append(lines, fmt.Sprintf("- %s: %d", key, report.ByExceptionKind[key]))
	}
	return strings.Join(lines, "\n") + "\n"
}

func Totals(reports []domain.DailyReport) (loaded, unloaded, exceptions int) {
	for _, report := range reports {
		loaded += report.Loaded
		unloaded += report.Unloaded
		exceptions += report.Exceptions
	}
	return loaded, unloaded, exceptions
}

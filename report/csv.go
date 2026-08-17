package report

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"transitmanifest/domain"
)

func WriteDaily(writer io.Writer, report domain.DailyReport) error {
	output := csv.NewWriter(writer)
	rows := [][]string{{"date", "created_manifests", "closed_manifests", "loaded", "unloaded", "handoffs", "exceptions"}, {
		report.Date, strconv.Itoa(report.CreatedManifests), strconv.Itoa(report.ClosedManifests), strconv.Itoa(report.Loaded), strconv.Itoa(report.Unloaded), strconv.Itoa(report.Handoffs), strconv.Itoa(report.Exceptions),
	}}
	for _, row := range rows {
		if err := output.Write(row); err != nil {
			return err
		}
	}
	output.Flush()
	if err := output.Error(); err != nil {
		return fmt.Errorf("write daily report: %w", err)
	}
	return nil
}

func ExceptionRows(report domain.DailyReport) [][2]string {
	result := make([][2]string, 0, len(report.ByExceptionKind))
	for kind, count := range report.ByExceptionKind {
		result = append(result, [2]string{kind, strconv.Itoa(count)})
	}
	return result
}

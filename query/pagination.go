package query

import (
	"fmt"
	"strconv"

	"transitmanifest/domain"
)

type Page struct {
	Items      []domain.Manifest `json:"items"`
	Offset     int               `json:"offset"`
	Limit      int               `json:"limit"`
	Total      int               `json:"total"`
	NextOffset *int              `json:"next_offset,omitempty"`
}

func Paginate(items []domain.Manifest, offset, limit int) (Page, error) {
	if offset < 0 {
		return Page{}, fmt.Errorf("offset %d must not be negative", offset)
	}
	if limit <= 0 || limit > 1000 {
		return Page{}, fmt.Errorf("limit %d must be between 1 and 1000", limit)
	}
	if offset > len(items) {
		offset = len(items)
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	page := Page{Items: append([]domain.Manifest(nil), items[offset:end]...), Offset: offset, Limit: limit, Total: len(items)}
	if end < len(items) {
		next := end
		page.NextOffset = &next
	}
	return page, nil
}

func ParsePage(value string, defaultValue int) (int, error) {
	if value == "" {
		return defaultValue, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return 0, fmt.Errorf("invalid page value %q", value)
	}
	return parsed, nil
}

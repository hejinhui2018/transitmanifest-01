package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"transitmanifest/domain"
)

type Record struct {
	Sequence     uint64          `json:"sequence"`
	WrittenAt    time.Time       `json:"written_at"`
	Type         string          `json:"type"`
	AggregateID  string          `json:"aggregate_id"`
	Actor        string          `json:"actor"`
	OccurredAt   time.Time       `json:"occurred_at"`
	Data         json.RawMessage `json:"data"`
	PreviousHash string          `json:"previous_hash"`
	Checksum     string          `json:"checksum"`
}

func newRecord(sequence uint64, previous string, event domain.Event, writtenAt time.Time) Record {
	return Record{
		Sequence:     sequence,
		WrittenAt:    writtenAt.UTC(),
		Type:         event.Type,
		AggregateID:  event.AggregateID,
		Actor:        event.Actor,
		OccurredAt:   event.OccurredAt.UTC(),
		Data:         append(json.RawMessage(nil), event.Data...),
		PreviousHash: previous,
	}
}

func (r Record) Event() domain.Event {
	return domain.Event{
		Type:        r.Type,
		AggregateID: r.AggregateID,
		Actor:       r.Actor,
		OccurredAt:  r.OccurredAt,
		Data:        append(json.RawMessage(nil), r.Data...),
	}
}

func (r Record) calculateChecksum() string {
	var canonical strings.Builder
	canonical.WriteString(strconv.FormatUint(r.Sequence, 10))
	canonical.WriteByte('\n')
	canonical.WriteString(r.WrittenAt.UTC().Format(time.RFC3339Nano))
	canonical.WriteByte('\n')
	canonical.WriteString(r.Type)
	canonical.WriteByte('\n')
	canonical.WriteString(r.AggregateID)
	canonical.WriteByte('\n')
	canonical.WriteString(r.Actor)
	canonical.WriteByte('\n')
	canonical.WriteString(r.OccurredAt.UTC().Format(time.RFC3339Nano))
	canonical.WriteByte('\n')
	canonical.Write(r.Data)
	canonical.WriteByte('\n')
	canonical.WriteString(r.PreviousHash)
	digest := sha256.Sum256([]byte(canonical.String()))
	return hex.EncodeToString(digest[:])
}

func (r Record) validate(expectedSequence uint64, previous string) error {
	if r.Sequence != expectedSequence {
		return fmt.Errorf("%w: sequence %d, expected %d", domain.ErrCorruptLog, r.Sequence, expectedSequence)
	}
	if r.PreviousHash != previous {
		return fmt.Errorf("%w: record %d previous checksum mismatch", domain.ErrCorruptLog, r.Sequence)
	}
	if err := domain.ValidateType(r.Type); err != nil {
		return err
	}
	want := r.calculateChecksum()
	if r.Checksum != want {
		return fmt.Errorf("%w: record %d checksum mismatch", domain.ErrCorruptLog, r.Sequence)
	}
	return nil
}

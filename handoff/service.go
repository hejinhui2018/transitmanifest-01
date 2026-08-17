package handoff

import (
	"time"

	"transitmanifest/domain"
	"transitmanifest/manifest"
)

type Service struct {
	manifests *manifest.Service
	now       func() time.Time
}

func New(manifests *manifest.Service) *Service {
	return &Service{manifests: manifests, now: time.Now}
}

func (s *Service) Sign(manifestID, from, to, signer, note string, at time.Time) (domain.HandoffReceipt, error) {
	if at.IsZero() {
		at = s.now()
	}
	receipt := domain.HandoffReceipt{ManifestID: manifestID, From: from, To: to, Signer: signer, Note: note, SignedAt: at}
	if err := validateReceipt(receipt); err != nil {
		return domain.HandoffReceipt{}, err
	}
	if err := s.manifests.SignHandoff(receipt); err != nil {
		return domain.HandoffReceipt{}, err
	}
	return receipt, nil
}

func validateReceipt(receipt domain.HandoffReceipt) error {
	for field, value := range map[string]string{"manifest_id": receipt.ManifestID, "from": receipt.From, "to": receipt.To, "signer": receipt.Signer} {
		if err := domain.ValidateIdentifier(field, value); err != nil {
			return err
		}
	}
	return nil
}

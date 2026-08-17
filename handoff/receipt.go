package handoff

import (
	"fmt"
	"strings"

	"transitmanifest/domain"
)

func Verify(receipt domain.HandoffReceipt, expectedFrom, expectedTo string) error {
	if receipt.From != expectedFrom || receipt.To != expectedTo {
		return fmt.Errorf("handoff path %s -> %s does not match expected %s -> %s", receipt.From, receipt.To, expectedFrom, expectedTo)
	}
	if strings.TrimSpace(receipt.Signer) == "" {
		return fmt.Errorf("handoff signer is required")
	}
	if receipt.SignedAt.IsZero() {
		return fmt.Errorf("handoff timestamp is required")
	}
	return nil
}

func Summary(receipt *domain.HandoffReceipt) string {
	if receipt == nil {
		return "unsigned"
	}
	return fmt.Sprintf("%s -> %s by %s at %s", receipt.From, receipt.To, receipt.Signer, receipt.SignedAt.Format("2006-01-02 15:04:05"))
}

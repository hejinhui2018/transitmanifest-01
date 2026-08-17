package handoff

import (
	"testing"
	"time"

	"transitmanifest/manifest"
	"transitmanifest/storage"
)

func TestSignAndVerify(t *testing.T) {
	store, _ := storage.Open(t.TempDir())
	defer store.Close()
	manifests, _ := manifest.Open(store)
	manifests.Create("m-1", "t-1", "plate-1", "a", "b", "operator", time.Now())
	service := New(manifests)
	receipt, err := service.Sign("m-1", "a", "b", "receiver", "sealed", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if err := Verify(receipt, "a", "b"); err != nil {
		t.Fatal(err)
	}
	if Summary(&receipt) == "unsigned" {
		t.Fatal("signed receipt has unsigned summary")
	}
}

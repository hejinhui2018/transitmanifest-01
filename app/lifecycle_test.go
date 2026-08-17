package app

import (
	"testing"
	"time"
)

func TestApplicationCloseAndReopen(t *testing.T) {
	directory := t.TempDir()
	config := Config{DataDir: directory, Actor: "tester"}
	application, err := Open(config)
	if err != nil {
		t.Fatal(err)
	}
	application.Manifests.Create("m-1", "trip-1", "plate-1", "a", "b", "tester", time.Now())
	if err := application.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(config)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	if _, err := reopened.Manifests.Get("m-1"); err != nil {
		t.Fatal(err)
	}
}

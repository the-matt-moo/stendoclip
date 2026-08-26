package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPersistenceRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "history.json")
	s := New()
	s.Push("history")
	s.Push("pin")
	s.TogglePin(0)
	if err := s.Save(path); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Len() != 2 || got.Get(0).Text != "pin" || !got.Get(0).Pinned || got.Get(1).Text != "history" {
		t.Fatalf("loaded stack: %#v %#v", got.Get(0), got.Get(1))
	}
}

func TestDedupAfterPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	s := New()
	s.Push("one")
	s.Push("two")
	if err := s.Save(path); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	got.Push("one")
	if got.Len() != 2 || got.Get(0).Text != "one" {
		t.Fatal("deduplication did not survive persistence")
	}
}

func TestLoadEmptyFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil || got.Len() != 0 {
		t.Fatalf("Load empty = len %d, err %v", got.Len(), err)
	}
}

func TestLoadCorruptJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "history.json")
	if err := os.WriteFile(path, []byte("{"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("corrupt JSON accepted")
	}
}

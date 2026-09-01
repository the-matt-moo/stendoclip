package store

import (
	"fmt"
	"strings"
	"testing"
)

func TestPushAndDedupMoveToTop(t *testing.T) {
	s := New()
	if !s.Push("first") || !s.Push("second") {
		t.Fatal("Push rejected valid text")
	}
	if got := s.Get(0).Text; got != "second" {
		t.Fatalf("top = %q, want second", got)
	}
	if !s.Push("first") || s.Len() != 2 || s.Get(0).Text != "first" {
		t.Fatal("duplicate was not moved to top")
	}
}

func TestPinUnpin(t *testing.T) {
	s := New()
	s.Push("one")
	s.TogglePin(0)
	if got := s.PinnedEntries(); len(got) != 1 || !got[0].Pinned || got[0].Text != "one" {
		t.Fatalf("pinned entries = %#v", got)
	}
	s.TogglePin(0)
	if len(s.PinnedEntries()) != 0 || s.Get(0).Pinned {
		t.Fatal("entry remained pinned")
	}
}

func TestDelete(t *testing.T) {
	s := New()
	s.Push("one")
	s.Push("two")
	s.Delete(0)
	if s.Len() != 1 || s.Get(0).Text != "one" {
		t.Fatalf("stack after delete = len %d, top %q", s.Len(), s.Get(0).Text)
	}
}

func TestCycleWraparound(t *testing.T) {
	s := New()
	for _, text := range []string{"one", "two", "three"} {
		s.Push(text)
	}
	if got := s.Cycle(2, 1); got != 0 {
		t.Fatalf("forward cycle = %d, want 0", got)
	}
	if got := s.Cycle(0, -1); got != 2 {
		t.Fatalf("backward cycle = %d, want 2", got)
	}
	if got := New().Cycle(0, 1); got != -1 {
		t.Fatalf("empty cycle = %d, want -1", got)
	}
}

func TestDefaultCapEvictsOldest(t *testing.T) {
	s := New()
	for i := 0; i < 51; i++ {
		s.Push(fmt.Sprintf("clip-%d", i))
	}
	if s.Len() != 50 {
		t.Fatalf("len = %d, want 50", s.Len())
	}
	if s.Get(49).Text == "clip-0" {
		t.Fatal("oldest entry was not evicted")
	}
}

func TestRejectsBlankOrOversizedClips(t *testing.T) {
	s := New()
	for _, text := range []string{"", " \n\t "} {
		if s.Push(text) {
			t.Fatalf("blank entry %q accepted", text)
		}
	}
	if s.Push(strings.Repeat("x", 65537)) {
		t.Fatal("oversized entry accepted")
	}
	if s.Len() != 0 {
		t.Fatalf("len = %d after rejected pushes", s.Len())
	}
}

func TestEntrySnapshotsAreCopies(t *testing.T) {
	s := New()
	s.Push("history")
	history := s.HistoryEntries()
	history[0].Text = "changed"
	if s.Get(0).Text != "history" {
		t.Fatal("history snapshot mutated stack")
	}
}

func TestClearHistoryPreservesPins(t *testing.T) {
	s := New()
	s.Push("history")
	s.Push("pin")
	s.TogglePin(0)
	s.ClearHistory()
	if s.Len() != 1 || s.Get(0).Text != "pin" || !s.Get(0).Pinned {
		t.Fatalf("clear result = %#v", s.Get(0))
	}
}

func TestSetMaxTrims(t *testing.T) {
	s := New()
	for _, text := range []string{"one", "two", "three"} {
		s.Push(text)
	}
	s.SetMax(2)
	if s.Len() != 2 || s.Get(0).Text != "three" || s.Get(1).Text != "two" {
		t.Fatalf("trimmed stack = %#v, %#v", s.Get(0), s.Get(1))
	}
}

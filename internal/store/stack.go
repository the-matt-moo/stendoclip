package store

import (
	"strings"
	"sync"
	"time"
)

const (
	defaultMaxSize  = 50
	defaultMaxBytes = 65536
)

type Entry struct {
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	Pinned    bool      `json:"pinned"`
}

type ClippingStack struct {
	mu       sync.RWMutex
	history  []Entry
	pins     []Entry
	maxSize  int
	maxBytes int
}

func New() *ClippingStack {
	return &ClippingStack{maxSize: defaultMaxSize, maxBytes: defaultMaxBytes}
}

func (s *ClippingStack) initDefaults() {
	if s.maxSize == 0 {
		s.maxSize = defaultMaxSize
	}
	if s.maxBytes == 0 {
		s.maxBytes = defaultMaxBytes
	}
}

func (s *ClippingStack) Push(text string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.initDefaults()
	if isBlankText(text) || len([]byte(text)) > s.maxBytes {
		return false
	}

	now := time.Now()
	for i, entry := range s.pins {
		if entry.Text == text {
			entry.Timestamp = now
			s.pins = append([]Entry{entry}, append(s.pins[:i], s.pins[i+1:]...)...)
			return true
		}
	}
	for i, entry := range s.history {
		if entry.Text == text {
			entry.Timestamp = now
			s.history = append([]Entry{entry}, append(s.history[:i], s.history[i+1:]...)...)
			return true
		}
	}

	s.history = append([]Entry{{Text: text, Timestamp: now}}, s.history...)
	s.trimLocked()
	return true
}

func (s *ClippingStack) Get(index int) Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if index < 0 {
		return Entry{}
	}
	if index < len(s.pins) {
		return s.pins[index]
	}
	index -= len(s.pins)
	if index >= len(s.history) {
		return Entry{}
	}
	return s.history[index]
}

func (s *ClippingStack) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.pins) + len(s.history)
}

func (s *ClippingStack) Delete(index int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 {
		return
	}
	if index < len(s.pins) {
		s.pins = append(s.pins[:index], s.pins[index+1:]...)
		return
	}
	index -= len(s.pins)
	if index < len(s.history) {
		s.history = append(s.history[:index], s.history[index+1:]...)
	}
}

func (s *ClippingStack) TogglePin(index int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.initDefaults()
	if index < 0 {
		return
	}
	if index < len(s.pins) {
		entry := s.pins[index]
		s.pins = append(s.pins[:index], s.pins[index+1:]...)
		entry.Pinned = false
		s.history = append([]Entry{entry}, s.history...)
		s.trimLocked()
		return
	}
	index -= len(s.pins)
	if index >= len(s.history) {
		return
	}
	entry := s.history[index]
	s.history = append(s.history[:index], s.history[index+1:]...)
	entry.Pinned = true
	s.pins = append([]Entry{entry}, s.pins...)
}

func (s *ClippingStack) PinnedEntries() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Entry(nil), s.pins...)
}

func (s *ClippingStack) HistoryEntries() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Entry(nil), s.history...)
}

func (s *ClippingStack) ClearHistory() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = nil
}

func (s *ClippingStack) Cycle(current, direction int) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := len(s.pins) + len(s.history)
	if n == 0 {
		return -1
	}
	return ((current+direction)%n + n) % n
}

func (s *ClippingStack) SetMax(n int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if n < 1 {
		n = 1
	}
	s.maxSize = n
	if s.maxBytes == 0 {
		s.maxBytes = defaultMaxBytes
	}
	s.trimLocked()
}

func (s *ClippingStack) SetMaxBytes(n int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if n < 1 {
		n = 1
	}
	s.maxBytes = n
}

func (s *ClippingStack) trimLocked() {
	if len(s.history) > s.maxSize {
		s.history = s.history[:s.maxSize]
	}
}

func isBlankText(text string) bool {
	return len(strings.TrimSpace(text)) == 0
}

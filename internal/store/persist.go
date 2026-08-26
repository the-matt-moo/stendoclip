package store

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type HistoryFile struct {
	History []Entry `json:"history"`
	Pins    []Entry `json:"pins"`
}

func (h HistoryFile) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}

func (s *ClippingStack) Save(path string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return (HistoryFile{
		History: append([]Entry(nil), s.history...),
		Pins:    append([]Entry(nil), s.pins...),
	}).Save(path)
}

func Load(path string) (*ClippingStack, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return New(), nil
	}
	if err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return New(), nil
	}

	var file HistoryFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, err
	}
	stack := New()
	stack.history = append([]Entry(nil), file.History...)
	stack.pins = append([]Entry(nil), file.Pins...)
	for i := range stack.history {
		stack.history[i].Pinned = false
	}
	for i := range stack.pins {
		stack.pins[i].Pinned = true
	}
	stack.trimLocked()
	return stack, nil
}

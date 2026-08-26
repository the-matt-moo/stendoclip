package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const maxLogSize = 1 << 20

type Logger struct {
	mu    sync.Mutex
	path  string
	file  *os.File
	debug bool
}

func New(path string) (*Logger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, err
	}
	return &Logger{path: path, file: file}, nil
}

func (l *Logger) SetDebug(enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.debug = enabled
}

func (l *Logger) Info(format string, args ...any) {
	l.write("INFO", format, args...)
}

func (l *Logger) Debug(format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.debug {
		l.writeLocked("DEBUG", format, args...)
	}
}

func (l *Logger) Error(format string, args ...any) {
	l.write("ERROR", format, args...)
}

func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}

func (l *Logger) write(level, format string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writeLocked(level, format, args...)
}

func (l *Logger) writeLocked(level, format string, args ...any) {
	if l.file == nil {
		return
	}
	message := fmt.Sprintf("%s [%s] %s\n", time.Now().Format(time.RFC3339), level, fmt.Sprintf(format, args...))
	if info, err := l.file.Stat(); err == nil && info.Size()+int64(len(message)) > maxLogSize {
		if err := l.rotateLocked(); err != nil {
			return
		}
	}
	_, _ = l.file.WriteString(message)
}

func (l *Logger) rotateLocked() error {
	if err := l.file.Close(); err != nil {
		return err
	}
	_ = os.Remove(l.path + ".1")
	if err := os.Rename(l.path, l.path+".1"); err != nil && !os.IsNotExist(err) {
		return err
	}
	file, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	l.file = file
	return nil
}

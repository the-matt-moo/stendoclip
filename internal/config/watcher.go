package config

import (
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher *fsnotify.Watcher
	stop    chan struct{}
	wg      sync.WaitGroup
}

func Watch(path string, debounce time.Duration, onChange func(Config), onError func(error)) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := fsw.Add(path); err != nil {
		_ = fsw.Close()
		return nil, err
	}
	w := &Watcher{watcher: fsw, stop: make(chan struct{})}
	w.wg.Add(1)
	go w.loop(path, debounce, onChange, onError)
	return w, nil
}

func (w *Watcher) Close() error {
	close(w.stop)
	err := w.watcher.Close()
	w.wg.Wait()
	return err
}

func (w *Watcher) loop(path string, debounce time.Duration, onChange func(Config), onError func(error)) {
	defer w.wg.Done()
	var timer *time.Timer
	for {
		select {
		case <-w.stop:
			if timer != nil {
				timer.Stop()
			}
			return
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
				continue
			}
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(debounce, func() {
				cfg, err := LoadConfig(path)
				if err != nil {
					onError(err)
					return
				}
				onChange(cfg)
			})
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			onError(err)
		}
	}
}

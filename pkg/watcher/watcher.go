package watcher

import (
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
	"regexp"
)

type Event struct {
	Path      string
	Operation string
}

type Watcher struct {
	fsWatcher    *fsnotify.Watcher
	Path         string
	IncludeRegex []string
	ExcludeRegex []string
}

func NewWatcher(path string, includeRegex, excludeRegex []string) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	dirs, err := FilePathWalkDir(path)
	if err != nil {
		return nil, err
	}

	for _, dir := range dirs {
		err = fsWatcher.Add(dir)
		if err != nil {
			return nil, err
		}
	}

	return &Watcher{
		fsWatcher:    fsWatcher,
		Path:         path,
		IncludeRegex: includeRegex,
		ExcludeRegex: excludeRegex,
	}, nil
}

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func (w *Watcher) Listen() (<-chan Event, error) {

	fsEvents := w.fsWatcher.Events
	return w.validateEvents(w.castEvents(fsEvents)), nil
}

func (w *Watcher) castEvents(fsEvents <-chan fsnotify.Event) <-chan Event {
	out := make(chan Event)
	go func() {
		for event := range fsEvents {
			out <- w.castEvent(event)
		}
		close(out)
	}()

	return out
}

func (w *Watcher) castEvent(event fsnotify.Event) Event {
	if event.Has(fsnotify.Create) {
		stat, _ := os.Stat(event.Name)
		if stat.IsDir() {
			_ = w.fsWatcher.Add(event.Name)
		}
	}
	return Event{
		Path:      event.Name,
		Operation: event.Op.String(),
	}
}

func (w *Watcher) validateEvents(events <-chan Event) <-chan Event {
	out := make(chan Event)
	go func() {
		for event := range events {
			flag := false
			for _, includeRegex := range w.IncludeRegex {
				ok, _ := regexp.Match(includeRegex, []byte(event.Path))
				if ok {
					flag = true
				}
			}
			if w.IncludeRegex == nil {
				flag = true
			}

			for _, excludeRegex := range w.ExcludeRegex {
				ok, _ := regexp.Match(excludeRegex, []byte(event.Path))
				if ok {
					flag = false
				}
			}

			if flag {
				out <- event
			}
		}

		close(out)
	}()

	return out
}

func (w *Watcher) Shutdown() error {
	return w.fsWatcher.Close()
}

// Package watch provides file-watching functionality for envoy-diff,
// allowing automatic re-diffing when environment files change on disk.
package watch

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Event represents a file change detected by the watcher.
type Event struct {
	// Path is the absolute path of the file that changed.
	Path string
	// ModTime is the modification time reported by the OS at detection time.
	ModTime time.Time
}

// Handler is called whenever one of the watched files changes.
type Handler func(event Event) error

// Watcher polls a set of file paths at a configurable interval and invokes
// a Handler when any file's modification time advances.
type Watcher struct {
	paths    []string
	interval time.Duration
	handler  Handler
	out      io.Writer
	stop     chan struct{}
}

// Options configures a new Watcher.
type Options struct {
	// Interval is how often the watcher polls the files. Defaults to 2s.
	Interval time.Duration
	// Out is the writer used for informational log lines. Defaults to os.Stderr.
	Out io.Writer
}

// New creates a Watcher that monitors the given paths and calls handler on
// each change. The watcher does not start until Run is called.
func New(paths []string, handler Handler, opts Options) (*Watcher, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("watch: at least one path is required")
	}
	if handler == nil {
		return nil, fmt.Errorf("watch: handler must not be nil")
	}

	interval := opts.Interval
	if interval <= 0 {
		interval = 2 * time.Second
	}

	out := opts.Out
	if out == nil {
		out = os.Stderr
	}

	// Resolve all paths to absolute form up front so comparisons are stable.
	abs := make([]string, len(paths))
	for i, p := range paths {
		a, err := filepath.Abs(p)
		if err != nil {
			return nil, fmt.Errorf("watch: resolving path %q: %w", p, err)
		}
		abs[i] = a
	}

	return &Watcher{
		paths:    abs,
		interval: interval,
		handler:  handler,
		out:      out,
		stop:     make(chan struct{}),
	}, nil
}

// Run starts the polling loop and blocks until Stop is called or a handler
// returns a non-nil error. It returns the handler error, or nil if stopped
// cleanly via Stop.
func (w *Watcher) Run() error {
	// Seed the initial modification times so we don't fire on first tick.
	last := make(map[string]time.Time, len(w.paths))
	for _, p := range w.paths {
		if info, err := os.Stat(p); err == nil {
			last[p] = info.ModTime()
		}
	}

	fmt.Fprintf(w.out, "watch: monitoring %d file(s) every %s\n", len(w.paths), w.interval)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-w.stop:
			fmt.Fprintln(w.out, "watch: stopped")
			return nil

		case <-ticker.C:
			for _, p := range w.paths {
				info, err := os.Stat(p)
				if err != nil {
					// File may have been temporarily unavailable; skip silently.
					continue
				}
				mod := info.ModTime()
				if prev, ok := last[p]; ok && !mod.After(prev) {
					continue
				}
				last[p] = mod
				fmt.Fprintf(w.out, "watch: change detected in %s\n", p)
				if err := w.handler(Event{Path: p, ModTime: mod}); err != nil {
					return fmt.Errorf("watch: handler error: %w", err)
				}
			}
		}
	}
}

// Stop signals the watcher to cease polling. It is safe to call from another
// goroutine. Calling Stop more than once is a no-op.
func (w *Watcher) Stop() {
	select {
	case <-w.stop:
		// already closed
	default:
		close(w.stop)
	}
}

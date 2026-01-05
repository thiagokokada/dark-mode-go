package dark

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestWatchDarkModeEmitsChanges(t *testing.T) {
	updates := make(chan struct{}, 10)
	errs := make(chan error, 1)
	stopped := make(chan struct{})

	factory := ChangeWatcherFactoryFunc(func() (ChangeWatcher, error) {
		return &mockChangeWatcher{
			updates: updates,
			errors:  errs,
			stop: func() error {
				close(stopped)
				return nil
			},
		}, nil
	})

	var mu sync.Mutex
	values := []bool{false, true, true}
	idx := 0
	checker := ModeCheckerFunc(func() (bool, error) {
		mu.Lock()
		defer mu.Unlock()
		if idx >= len(values) {
			return values[len(values)-1], nil
		}
		v := values[idx]
		idx++
		return v, nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watcher := NewWatcher(checker, factory)
	events, watchErrs, err := watcher.Watch(ctx)
	if err != nil {
		t.Fatalf("WatchDarkMode error: %v", err)
	}
	_ = watchErrs

	if got := readBool(t, events); got != false {
		t.Fatalf("initial event = %v, want false", got)
	}

	updates <- struct{}{}
	if got := readBool(t, events); got != true {
		t.Fatalf("updated event = %v, want true", got)
	}

	updates <- struct{}{}
	select {
	case got := <-events:
		t.Fatalf("unexpected duplicate event: %v", got)
	case <-time.After(50 * time.Millisecond):
	}

	cancel()
	waitClosed(t, stopped)
	waitClosedBool(t, events)
}

type mockChangeWatcher struct {
	updates <-chan struct{}
	errors  <-chan error
	stop    func() error
}

func (mw *mockChangeWatcher) Updates() <-chan struct{} {
	return mw.updates
}

func (mw *mockChangeWatcher) Errors() <-chan error {
	return mw.errors
}

func (mw *mockChangeWatcher) Stop() error {
	if mw.stop == nil {
		return nil
	}
	return mw.stop()
}

func readBool(t *testing.T, ch <-chan bool) bool {
	t.Helper()
	select {
	case v, ok := <-ch:
		if !ok {
			t.Fatal("events channel closed unexpectedly")
		}
		return v
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for event")
		return false
	}
}

func waitClosed(t *testing.T, ch <-chan struct{}) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for watcher stop")
	}
}

func waitClosedBool(t *testing.T, ch <-chan bool) {
	t.Helper()
	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("events channel still open")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for events channel close")
	}
}

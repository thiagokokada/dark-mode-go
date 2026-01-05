package dark

import "context"

type ModeChecker interface {
	IsDarkMode() (bool, error)
}

type ModeCheckerFunc func() (bool, error)

func (fn ModeCheckerFunc) IsDarkMode() (bool, error) {
	return fn()
}

type ChangeWatcher interface {
	Updates() <-chan struct{}
	Errors() <-chan error
	Stop() error
}

type ChangeWatcherFactory interface {
	Start() (ChangeWatcher, error)
}

type ChangeWatcherFactoryFunc func() (ChangeWatcher, error)

func (fn ChangeWatcherFactoryFunc) Start() (ChangeWatcher, error) {
	return fn()
}

type changeWatcher struct {
	updates <-chan struct{}
	errors  <-chan error
	stop    func() error
}

func (cw *changeWatcher) Updates() <-chan struct{} {
	return cw.updates
}

func (cw *changeWatcher) Errors() <-chan error {
	return cw.errors
}

func (cw *changeWatcher) Stop() error {
	if cw.stop == nil {
		return nil
	}
	return cw.stop()
}

type Watcher struct {
	checker ModeChecker
	factory ChangeWatcherFactory
}

func NewWatcher(checker ModeChecker, factory ChangeWatcherFactory) *Watcher {
	if checker == nil {
		checker = ModeCheckerFunc(IsDarkMode)
	}
	if factory == nil {
		factory = ChangeWatcherFactoryFunc(startChangeWatcher)
	}
	return &Watcher{
		checker: checker,
		factory: factory,
	}
}

// WatchDarkMode emits the current dark mode setting and any changes until ctx is done.
// The events channel sends the current state immediately after setup.
func WatchDarkMode(ctx context.Context) (<-chan bool, <-chan error, error) {
	return NewWatcher(nil, nil).Watch(ctx)
}

// Watch emits the current dark mode setting and any changes until ctx is done.
// The events channel sends the current state immediately after setup.
func (w *Watcher) Watch(ctx context.Context) (<-chan bool, <-chan error, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	watcher, err := w.factory.Start()
	if err != nil {
		return nil, nil, err
	}

	events := make(chan bool, 1)
	errs := make(chan error, 1)

	go func() {
		sendErr := func(err error) {
			if err == nil {
				return
			}
			select {
			case errs <- err:
			default:
			}
		}
		sendEvent := func(value bool) bool {
			select {
			case events <- value:
				return true
			case <-ctx.Done():
				return false
			}
		}

		defer func() {
			if err := watcher.Stop(); err != nil {
				sendErr(err)
			}
			close(events)
			close(errs)
		}()

		current, err := w.checker.IsDarkMode()
		if err != nil {
			sendErr(err)
			return
		}
		last := current
		if !sendEvent(current) {
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case err, ok := <-watcher.Errors():
				if ok {
					sendErr(err)
				}
				return
			case _, ok := <-watcher.Updates():
				if !ok {
					return
				}
				current, err := w.checker.IsDarkMode()
				if err != nil {
					sendErr(err)
					return
				}
				if current != last {
					last = current
					if !sendEvent(current) {
						return
					}
				}
			}
		}
	}()

	return events, errs, nil
}

package dark

import (
	"sync"
	"time"
)

const macOSPollInterval = time.Second

func startChangeWatcher() (ChangeWatcher, error) {
	updates := make(chan struct{}, 1)
	errs := make(chan error, 1)
	stopCh := make(chan struct{})
	var once sync.Once

	go func() {
		defer close(updates)
		defer close(errs)

		ticker := time.NewTicker(macOSPollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				select {
				case updates <- struct{}{}:
				default:
				}
			case <-stopCh:
				return
			}
		}
	}()

	stop := func() error {
		once.Do(func() {
			close(stopCh)
		})
		return nil
	}

	return &changeWatcher{
		updates: updates,
		errors:  errs,
		stop:    stop,
	}, nil
}

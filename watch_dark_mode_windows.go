package dark

import (
	"fmt"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

func startChangeWatcher() (ChangeWatcher, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, registryPath, registry.QUERY_VALUE|registry.NOTIFY)
	if err != nil {
		return nil, fmt.Errorf("%w: registry open key path '%s': %w", OsError, registryPath, err)
	}

	changeEvent, err := windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		_ = key.Close()
		return nil, fmt.Errorf("%w: registry change event: %w", OsError, err)
	}

	stopEvent, err := windows.CreateEvent(nil, 1, 0, nil)
	if err != nil {
		_ = windows.CloseHandle(changeEvent)
		_ = key.Close()
		return nil, fmt.Errorf("%w: registry stop event: %w", OsError, err)
	}

	updates := make(chan struct{}, 1)
	errs := make(chan error, 1)

	go func() {
		emitErr := func(err error) {
			if err == nil {
				return
			}
			select {
			case errs <- err:
			default:
			}
		}

		defer close(updates)
		defer close(errs)
		defer windows.CloseHandle(changeEvent)
		defer windows.CloseHandle(stopEvent)
		defer key.Close()

		if err := windows.RegNotifyChangeKeyValue(
			windows.Handle(key),
			false,
			windows.REG_NOTIFY_CHANGE_LAST_SET,
			changeEvent,
			true,
		); err != nil {
			emitErr(fmt.Errorf("%w: registry notify: %w", OsError, err))
			return
		}

		handles := []windows.Handle{changeEvent, stopEvent}
		for {
			status, err := windows.WaitForMultipleObjects(handles, false, windows.INFINITE)
			if err != nil {
				emitErr(fmt.Errorf("%w: registry wait: %w", OsError, err))
				return
			}
			switch status {
			case windows.WAIT_OBJECT_0:
				select {
				case updates <- struct{}{}:
				default:
				}
				if err := windows.RegNotifyChangeKeyValue(
					windows.Handle(key),
					false,
					windows.REG_NOTIFY_CHANGE_LAST_SET,
					changeEvent,
					true,
				); err != nil {
					emitErr(fmt.Errorf("%w: registry notify: %w", OsError, err))
					return
				}
			case windows.WAIT_OBJECT_0 + 1:
				return
			default:
				emitErr(fmt.Errorf("%w: registry wait status %d", OsError, status))
				return
			}
		}
	}()

	stop := func() error {
		return windows.SetEvent(stopEvent)
	}

	return &changeWatcher{
		updates: updates,
		errors:  errs,
		stop:    stop,
	}, nil
}

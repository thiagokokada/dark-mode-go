package dark

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

func startChangeWatcher() (ChangeWatcher, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("%w: dbus session bus connection: %w", OsError, err)
	}

	matchOptions := []dbus.MatchOption{
		dbus.WithMatchInterface(dbusSettingsInterface),
		dbus.WithMatchMember(dbusSettingsMember),
		dbus.WithMatchObjectPath(dbusPortalObjectPath),
	}
	if err := conn.AddMatchSignal(matchOptions...); err != nil {
		conn.Close()
		return nil, fmt.Errorf("%w: dbus match signal: %w", OsError, err)
	}

	signalCh := make(chan *dbus.Signal, 8)
	conn.Signal(signalCh)

	updates := make(chan struct{}, 1)
	errs := make(chan error, 1)

	go func() {
		defer close(updates)
		defer close(errs)

		for sig := range signalCh {
			if sig == nil || len(sig.Body) < 2 {
				continue
			}
			namespace, ok := sig.Body[0].(string)
			if !ok || namespace != dbusAppearanceNamespace {
				continue
			}
			key, ok := sig.Body[1].(string)
			if !ok || key != dbusColorSchemeKey {
				continue
			}
			select {
			case updates <- struct{}{}:
			default:
			}
		}
	}()

	stop := func() error {
		_ = conn.RemoveMatchSignal(matchOptions...)
		conn.Close()
		return nil
	}

	return &changeWatcher{
		updates: updates,
		errors:  errs,
		stop:    stop,
	}, nil
}

package dark

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

const dbusMethod = "org.freedesktop.portal.Settings.Read"

func IsDarkMode() (bool, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return false, fmt.Errorf("%w: dbus session bus connection: %w", OsError, err)
	}
	defer conn.Close()

	obj := conn.Object(
		"org.freedesktop.portal.Desktop",
		"/org/freedesktop/portal/desktop",
	)
	var colorScheme uint32
	err = obj.Call(dbusMethod, 0, "org.freedesktop.appearance", "color-scheme").Store(&colorScheme)
	if err != nil {
		return false, fmt.Errorf("%w: dbus method '%s' call: %w", OsError, dbusMethod, err)
	}
	// 0: no preference, 1: prefer dark mode, 2: prefer light mode
	return colorScheme == 0 || colorScheme == 1, nil
}

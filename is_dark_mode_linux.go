package dark

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

const (
	dbusSettingsInterface   = "org.freedesktop.portal.Settings"
	dbusSettingsMethod      = dbusSettingsInterface + ".Read"
	dbusSettingsMember      = "SettingChanged"
	dbusPortalObjectPath    = "/org/freedesktop/portal/desktop"
	dbusAppearanceNamespace = "org.freedesktop.appearance"
	dbusColorSchemeKey      = "color-scheme"
)

func IsDarkMode() (bool, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return false, fmt.Errorf("%w: dbus session bus connection: %w", OsError, err)
	}
	defer conn.Close()

	obj := conn.Object(dbusAppearanceNamespace, dbusPortalObjectPath)
	var colorScheme uint32
	err = obj.Call(dbusSettingsMethod, 0, dbusAppearanceNamespace, dbusColorSchemeKey).Store(&colorScheme)
	if err != nil {
		return false, fmt.Errorf("%w: dbus method '%s' call: %w", OsError, dbusSettingsMethod, err)
	}
	// 0: no preference, 1: prefer dark mode, 2: prefer light mode
	return colorScheme == 0 || colorScheme == 1, nil
}

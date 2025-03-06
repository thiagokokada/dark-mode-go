package dark

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

const registryPath = `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`

func IsDarkMode() (bool, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, registryPath, registry.QUERY_VALUE)
	if err != nil {
		return false, fmt.Errorf(
			"%w: registry open key path '%s': %w", OsError, registryPath, err,
		)
	}
	defer k.Close()

	v, _, err := k.GetIntegerValue("AppsUseLightTheme")
	if err != nil {
		return false, fmt.Errorf(
			"%w: registry get integer value: %w", OsError, err,
		)
	}

	// 0 means dark mode, 1 means light mode
	return v == 0, nil
}

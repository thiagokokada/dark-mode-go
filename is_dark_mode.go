//go:build !windows && !darwin && !linux

package dark

import (
	"fmt"
	"runtime"
)

// fallback if we have no way to detect dark-mode
func IsDarkMode() (bool, error) {
	return false, fmt.Errorf("%w: GOOS=%s", NotImplementedError, runtime.GOOS)
}

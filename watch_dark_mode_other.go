//go:build !darwin && !linux && !windows

package dark

import (
	"fmt"
	"runtime"
)

func startChangeWatcher() (ChangeWatcher, error) {
	return nil, fmt.Errorf("%w: GOOS=%s", NotImplementedError, runtime.GOOS)
}

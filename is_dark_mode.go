//go:build !windows && !darwin && !linux
package dark

// fallback if we have no way to detect dark
func IsDarkMode() (bool, error) {
	return false, NotImplementedError
}

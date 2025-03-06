package dark_test

import (
	"testing"

	"github.com/thiagokokada/dark-mode-go"
)

func TestIsDarkMode(t *testing.T) {
	r, err := dark.IsDarkMode()
	t.Logf("Result: %v", r)
	t.Logf("Error: %v", err)
}

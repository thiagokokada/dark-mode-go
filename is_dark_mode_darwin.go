package dark

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

func IsDarkMode() (bool, error) {
	cmd := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			return false, err
		}
		return false, nil
	}
	return strings.TrimSpace(out.String()) == "Dark", nil
}

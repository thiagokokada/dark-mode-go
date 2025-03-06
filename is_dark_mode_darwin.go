package dark

import (
	"bytes"
	"os/exec"
	"strings"
)

func IsDarkMode() (bool, error) {
	cmd := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	return err == nil && strings.TrimSpace(out.String()) == "Dark", nil
}

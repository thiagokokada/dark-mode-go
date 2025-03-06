# dark-mode-go

This is a small library for Go to detect dark mode in different operating
systems.

## Supported OSes

- [Linux](./is_dark_mode_linux.go)
- [macOS](./is_dark_mode_darwin.go)
- [Windows](./is_dark_mode_windows.go)

## Usage

```go
package main

import (
	"fmt"

	"github.com/thiagokokada/dark-mode-go"
)

func main() {
    r, err := dark.IsDarkMode()
    if err != nil {
        panic(err)
    }
    if r {
        fmt.Println("Dark mode")
    } else {
        fmt.Println("Light mode")
    }
}
```

## TODO

- [ ] Support reacting to theme change events
- [ ] Support more operating systems

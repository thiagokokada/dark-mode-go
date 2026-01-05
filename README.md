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
	"context"
	"fmt"
	"log"

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

## Watching for theme changes

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

events, errs, err := dark.WatchDarkMode(ctx)
if err != nil {
	log.Fatal(err)
}

for {
	select {
	case isDark, ok := <-events:
		if !ok {
			return
		}
		fmt.Println("Dark mode:", isDark)
	case err, ok := <-errs:
		if ok && err != nil {
			log.Printf("watch error: %v", err)
		}
	}
}
```

## TODO

- [ ] Support more operating systems

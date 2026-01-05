package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/thiagokokada/dark-mode-go"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signals
		cancel()
	}()

	events, errs, err := dark.WatchDarkMode(ctx)
	if err != nil {
		log.Fatalf("watch error: %v", err)
	}

	for {
		select {
		case isDark, ok := <-events:
			if !ok {
				return
			}
			fmt.Printf("dark mode: %v\n", isDark)
		case err, ok := <-errs:
			if !ok {
				return
			}
			if err != nil {
				log.Printf("watch error: %v", err)
			}
		}
	}
}

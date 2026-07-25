package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ndyakov/lastfm-mcp/internal/app"
)

var version = "dev"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, os.Args[1:], os.Getenv, version); err != nil {
		fmt.Fprintln(os.Stderr, "lastfm-mcp:", err)
		os.Exit(1)
	}
}

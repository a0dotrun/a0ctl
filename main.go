package main

import (
	"context"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/a0dotrun/a0ctl/internal/cli"
)

func main() {
	os.Exit(run())
}

func run() (exitCode int) {
	ctx, cancel := newContext()
	defer cancel()

	return cli.Run(ctx, os.Args[1:]...)
}

func newContext() (context.Context, context.CancelFunc) {
	// NOTE: when signal.Notify is called for os.Interrupt it traps both
	// ^C (Control-C) and ^BREAK (Control-Break) on Windows.
	signals := []os.Signal{os.Interrupt}
	if runtime.GOOS != "windows" {
		signals = append(signals, syscall.SIGTERM)
	}

	return signal.NotifyContext(context.Background(), signals...)
}

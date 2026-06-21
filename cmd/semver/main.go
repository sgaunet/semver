// Command semver manipulates semantic versions per https://semver.org: bump,
// pre-release lifecycle, compare, sort, validate, get, and constraint testing.
//
// stdout carries machine-parseable data; stderr carries diagnostics. Exit codes are
// documented in 'semver --help'.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sgaunet/semver/internal/cli"
)

func main() {
	os.Exit(run())
}

// run wires signal handling and dispatches to the CLI, returning the exit code. It
// is separate from main so that deferred cleanup (signal stop) runs before os.Exit.
func run() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return cli.Run(ctx, os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
}

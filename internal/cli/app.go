// Package cli implements the semver command-line interface. It is a thin wrapper
// over the internal/version and internal/constraint packages; all domain logic
// lives there. This package owns flag parsing, the stdout/stderr split, output
// formatting, and exit codes.
package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
)

// Build information, overridable at link time via -ldflags "-X .../cli.Version=...".
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

type commandFunc func(ctx context.Context, args []string, s streams) int

type command struct {
	name    string
	summary string
	run     commandFunc
}

func commands() []command {
	return []command{
		{"major", "Bump the major version", runBump("major")},
		{"minor", "Bump the minor version", runBump("minor")},
		{"patch", "Bump the patch version", runBump("patch")},
		{"prerelease", "Start or increment a pre-release", runPrerelease},
		{"release", "Finalize a pre-release to its stable version", runRelease},
		{"compare", "Compare two versions", runCompare},
		{"sort", "Sort versions by precedence", runSort},
		{"validate", "Validate and parse a version", runValidate},
		{"get", "Extract a component of a version", runGet},
		{"satisfies", "Test a version against a constraint", runSatisfies},
		{"version", "Print build version information", runVersion},
	}
}

// Run is the entry point: it dispatches args[0] to a subcommand and returns the
// process exit code. in/out/err are stdin/stdout/stderr.
func Run(ctx context.Context, args []string, in io.Reader, out, errw io.Writer) int {
	if len(args) == 0 {
		printMainHelp(errw)
		return CodeUsage
	}
	switch args[0] {
	case "-h", "--help", "help":
		printMainHelp(out)
		return CodeOK
	case "--version":
		s := streams{in: in, out: out, err: errw, format: "text"}
		return runVersion(ctx, nil, s)
	}
	for _, c := range commands() {
		if c.name == args[0] {
			s := streams{in: in, out: out, err: errw, noColor: noColorFromEnv()}
			return c.run(ctx, args[1:], s)
		}
	}
	fmt.Fprintf(errw, "semver: unknown command %q\nRun 'semver --help' for usage.\n", args[0])
	return CodeUsage
}

// setup builds a FlagSet for a subcommand with the global flags bound to s and a
// help printer that documents the command.
func setup(name, argsHint, summary string, s *streams) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.StringVar(&s.format, "output", "text", "output format: text|json")
	fs.BoolVar(&s.quiet, "quiet", false, "suppress non-error stderr output")
	fs.BoolVar(&s.quiet, "q", false, "shorthand for --quiet")
	fs.BoolVar(&s.verbose, "verbose", false, "verbose stderr output")
	fs.BoolVar(&s.verbose, "v", false, "shorthand for --verbose")
	fs.Usage = func() {
		w := fs.Output()
		fmt.Fprintf(w, "%s - %s\n\n", name, summary)
		fmt.Fprintf(w, "Usage:\n  semver %s [flags] %s\n\nFlags:\n", name, argsHint)
		fs.PrintDefaults()
		fmt.Fprintf(w, "\nExit codes: see 'semver --help'.\n")
	}
	return fs
}

// parse parses args, routing help to stdout and errors to stderr, then validates the
// global flags. Flags may appear before, between, or after positional arguments. It
// returns the positional arguments and (code, ok); when ok is false, code is the
// value to return.
func parse(fs *flag.FlagSet, args []string, s *streams) ([]string, int, bool) {
	if hasHelp(args) {
		fs.SetOutput(s.out)
		fs.Usage()
		return nil, CodeOK, false
	}
	fs.SetOutput(s.err)

	// Interleave flags and positionals: stdlib flag stops at the first non-flag, so
	// re-parse the remainder after each positional to allow trailing flags.
	var positional []string
	rest := args
	for {
		if err := fs.Parse(rest); err != nil {
			return nil, CodeUsage, false
		}
		rest = fs.Args()
		if len(rest) == 0 {
			break
		}
		positional = append(positional, rest[0])
		rest = rest[1:]
	}

	if s.format != "text" && s.format != "json" {
		s.errorf("invalid --output %q: want text or json", s.format)
		return nil, CodeUsage, false
	}
	if s.quiet && s.verbose {
		s.errorf("--quiet and --verbose are mutually exclusive")
		return nil, CodeUsage, false
	}
	return positional, CodeOK, true
}

func hasHelp(args []string) bool {
	for _, a := range args {
		if a == "--" {
			break
		}
		if a == "-h" || a == "--help" {
			return true
		}
	}
	return false
}

func printMainHelp(w io.Writer) {
	fmt.Fprint(w, mainHelpText)
}

const mainHelpText = `semver - manipulate semantic versions (https://semver.org)

Usage:
  semver <command> [flags] [args...]

Commands:
  major <version>              Bump the major version
  minor <version>              Bump the minor version
  patch <version>              Bump the patch version
  prerelease <version> [--pre id] [--bump major|minor|patch]
                               Start (--pre) or increment a pre-release
  release <version>            Finalize a pre-release to its stable version
  compare <a> <b>              Compare two versions
  sort [versions...]           Sort versions by precedence (reads stdin if no args)
  validate <version>           Validate and parse a version
  get <component> <version>    Extract major|minor|patch|prerelease|build
  satisfies <version> <range>  Test a version against a constraint range
  version                      Print build version information

Global flags:
  --output text|json   Output format (default text)
  --quiet, -q          Suppress non-error stderr output
  --verbose, -v        Verbose stderr output
  --help, -h           Show help

Exit codes:
  0   success / equal / valid / satisfied / present
  1   generic failure
  2   usage error (bad flags or args, malformed constraint, invalid version)
  10  compare: lower | satisfies: no | validate: invalid | get: absent
  11  compare: higher

Configuration precedence: flags > environment (NO_COLOR) > defaults.
`

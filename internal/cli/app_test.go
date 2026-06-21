package cli_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/sgaunet/semver/internal/cli"
)

// run invokes the CLI with the given args and no stdin.
func run(args ...string) (stdout, stderr string, code int) {
	return runStdin("", args...)
}

// runStdin invokes the CLI with the given stdin and args, returning captured
// stdout, stderr, and the exit code.
func runStdin(stdin string, args ...string) (string, string, int) {
	var out, errb bytes.Buffer
	code := cli.Run(context.Background(), args, strings.NewReader(stdin), &out, &errb)
	return out.String(), errb.String(), code
}

func TestNoArgsIsUsageError(t *testing.T) {
	out, errb, code := run()
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
	if out != "" {
		t.Errorf("stdout = %q, want empty", out)
	}
	if errb == "" {
		t.Error("want help text on stderr")
	}
}

func TestUnknownCommand(t *testing.T) {
	out, errb, code := run("frobnicate")
	if code != 2 || out != "" || !strings.Contains(errb, "unknown command") {
		t.Errorf("got out=%q err=%q code=%d", out, errb, code)
	}
}

func TestTopLevelHelpToStdout(t *testing.T) {
	out, _, code := run("--help")
	if code != 0 || !strings.Contains(out, "Usage:") || !strings.Contains(out, "Exit codes:") {
		t.Errorf("help: code=%d out=%q", code, out)
	}
}

func TestInvalidOutputFormat(t *testing.T) {
	_, _, code := run("patch", "v1.0.0", "--output", "yaml")
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}

func TestQuietVerboseMutuallyExclusive(t *testing.T) {
	_, _, code := run("patch", "v1.0.0", "-q", "-v")
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}

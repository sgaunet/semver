package integration

import (
	"strings"
	"testing"
)

// TestIntegrationHelpForEachCommand asserts every command's --help exits 0 and
// writes usage to stdout (FR-018 / SC-008).
func TestIntegrationHelpForEachCommand(t *testing.T) {
	cmds := []string{
		"major", "minor", "patch", "prerelease", "release",
		"compare", "sort", "validate", "get", "satisfies", "version",
	}
	for _, c := range cmds {
		out, _, code := run(t, "", c, "--help")
		if code != 0 {
			t.Errorf("%s --help: exit = %d, want 0", c, code)
		}
		if !strings.Contains(out, "Usage:") {
			t.Errorf("%s --help: stdout missing Usage section: %q", c, out)
		}
	}
}

func TestIntegrationTopLevelHelpAndVersion(t *testing.T) {
	out, _, code := run(t, "", "--help")
	if code != 0 || !strings.Contains(out, "Exit codes:") {
		t.Errorf("--help: code=%d out=%q", code, out)
	}
	out, _, code = run(t, "", "--version")
	if code != 0 || !strings.Contains(out, "semver") {
		t.Errorf("--version: code=%d out=%q", code, out)
	}
}

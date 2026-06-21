package integration

import "testing"

// TestIntegrationLifecycle walks the full pre-release lifecycle through the binary
// and confirms it round-trips to the anchored stable version (spec US3 / SC-005).
func TestIntegrationLifecycle(t *testing.T) {
	start, _, code := run(t, "", "prerelease", "v1.0.0", "--pre", "rc", "--bump", "minor")
	if code != 0 || start != "v1.1.0-rc.1\n" {
		t.Fatalf("start: out=%q code=%d, want v1.1.0-rc.1", start, code)
	}
	inc, _, code := run(t, "", "prerelease", "v1.1.0-rc.1")
	if code != 0 || inc != "v1.1.0-rc.2\n" {
		t.Fatalf("increment: out=%q code=%d, want v1.1.0-rc.2", inc, code)
	}
	final, _, code := run(t, "", "release", "v1.1.0-rc.2")
	if code != 0 || final != "v1.1.0\n" {
		t.Fatalf("finalize: out=%q code=%d, want v1.1.0", final, code)
	}
}

func TestIntegrationPrereleaseBetaStart(t *testing.T) {
	out, _, code := run(t, "", "prerelease", "v1.0.0", "--pre", "beta")
	if code != 0 || out != "v1.0.1-beta.1\n" {
		t.Errorf("beta start: out=%q code=%d, want v1.0.1-beta.1", out, code)
	}
}

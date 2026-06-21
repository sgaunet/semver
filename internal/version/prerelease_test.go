package version_test

import (
	"testing"

	"github.com/sgaunet/semver/internal/version"
)

func TestStartPrerelease(t *testing.T) {
	cases := []struct {
		in   string
		id   string
		kind version.BumpKind
		want string
	}{
		{"v1.0.0", "rc", version.BumpMinor, "v1.1.0-rc.1"},
		{"v1.0.0", "beta", version.BumpPatch, "v1.0.1-beta.1"},
		{"1.2.3", "alpha", version.BumpMajor, "2.0.0-alpha.1"},
	}
	for _, c := range cases {
		got, err := version.MustParse(c.in).StartPrerelease(c.id, c.kind)
		if err != nil {
			t.Errorf("StartPrerelease(%q,%q): %v", c.in, c.id, err)
			continue
		}
		if got.String() != c.want {
			t.Errorf("StartPrerelease(%q,%q) = %q, want %q", c.in, c.id, got.String(), c.want)
		}
	}
}

func TestStartPrereleaseInvalidID(t *testing.T) {
	for _, id := range []string{"", "rc!", "01"} {
		if _, err := version.MustParse("1.0.0").StartPrerelease(id, version.BumpPatch); err == nil {
			t.Errorf("StartPrerelease with id %q = nil error, want error", id)
		}
	}
}

func TestIncPrerelease(t *testing.T) {
	cases := []struct{ in, want string }{
		{"v1.1.0-rc.1", "v1.1.0-rc.2"},
		{"1.0.0-beta", "1.0.0-beta.1"},
		{"1.0.0-alpha.9", "1.0.0-alpha.10"},
	}
	for _, c := range cases {
		got, err := version.MustParse(c.in).IncPrerelease()
		if err != nil {
			t.Errorf("IncPrerelease(%q): %v", c.in, err)
			continue
		}
		if got.String() != c.want {
			t.Errorf("IncPrerelease(%q) = %q, want %q", c.in, got.String(), c.want)
		}
	}
}

func TestIncPrereleaseOnStableErrors(t *testing.T) {
	if _, err := version.MustParse("1.0.0").IncPrerelease(); err == nil {
		t.Error("IncPrerelease on stable version = nil error, want error")
	}
}

func TestFinalize(t *testing.T) {
	got, err := version.MustParse("v1.1.0-rc.2").Finalize()
	if err != nil {
		t.Fatalf("Finalize: %v", err)
	}
	if got.String() != "v1.1.0" {
		t.Errorf("Finalize = %q, want v1.1.0", got.String())
	}
}

func TestFinalizeOnStableErrors(t *testing.T) {
	if _, err := version.MustParse("1.0.0").Finalize(); err == nil {
		t.Error("Finalize on stable version = nil error, want error")
	}
}

// TestLifecycleRoundTrip verifies start -> increment -> finalize returns the anchored
// stable version (spec SC-005).
func TestLifecycleRoundTrip(t *testing.T) {
	start, err := version.MustParse("v1.0.0").StartPrerelease("rc", version.BumpMinor)
	if err != nil {
		t.Fatal(err)
	}
	inc, err := start.IncPrerelease()
	if err != nil {
		t.Fatal(err)
	}
	final, err := inc.Finalize()
	if err != nil {
		t.Fatal(err)
	}
	if final.String() != "v1.1.0" {
		t.Errorf("lifecycle final = %q, want v1.1.0", final.String())
	}
}

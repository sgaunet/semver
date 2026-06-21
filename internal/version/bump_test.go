package version_test

import (
	"testing"

	"github.com/sgaunet/semver/internal/version"
)

func TestBumps(t *testing.T) {
	cases := []struct {
		in   string
		kind string
		want string
	}{
		// Stable bumps reset lower components and preserve the prefix.
		{"v1.0.0", "patch", "v1.0.1"},
		{"v1.2.3", "minor", "v1.3.0"},
		{"v1.2.3", "major", "v2.0.0"},
		{"1.0.0", "patch", "1.0.1"}, // no prefix in -> none out
		// Pre-release collapse (node-semver convention).
		{"v1.2.0-rc.1", "patch", "v1.2.0"},
		{"v1.2.3-rc.1", "patch", "v1.2.3"},
		{"v1.2.0-rc.1", "minor", "v1.2.0"},
		{"v1.2.1-rc.1", "minor", "v1.3.0"},
		{"v1.0.0-rc.1", "major", "v1.0.0"},
		{"v1.2.0-rc.1", "major", "v2.0.0"},
		// Build metadata is dropped by any bump.
		{"1.2.3+build.7", "patch", "1.2.4"},
	}
	for _, c := range cases {
		v := version.MustParse(c.in)
		var got version.Version
		switch c.kind {
		case "major":
			got = v.IncMajor()
		case "minor":
			got = v.IncMinor()
		default:
			got = v.IncPatch()
		}
		if got.String() != c.want {
			t.Errorf("%s %s = %q, want %q", c.kind, c.in, got.String(), c.want)
		}
	}
}

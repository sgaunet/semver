package version_test

import (
	"testing"

	"github.com/sgaunet/semver/internal/version"
)

func TestParseValidAndRoundTrip(t *testing.T) {
	cases := []string{
		"0.0.0",
		"1.2.3",
		"v1.2.3",
		"V1.2.3",
		"1.0.0-rc.1",
		"v1.2.3-rc.1+build.7",
		"1.0.0-alpha.beta-1",
		"1.0.0+20130313144700",
		"10.20.30",
	}
	for _, in := range cases {
		v, err := version.Parse(in)
		if err != nil {
			t.Errorf("Parse(%q) unexpected error: %v", in, err)
			continue
		}
		if got := v.String(); got != in {
			t.Errorf("round-trip: Parse(%q).String() = %q, want %q", in, got, in)
		}
	}
}

func TestParseComponents(t *testing.T) {
	v, err := version.Parse("v1.2.3-rc.1+build.7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Major != 1 || v.Minor != 2 || v.Patch != 3 {
		t.Errorf("core = %d.%d.%d, want 1.2.3", v.Major, v.Minor, v.Patch)
	}
	if v.Prerelease() != "rc.1" {
		t.Errorf("prerelease = %q, want rc.1", v.Prerelease())
	}
	if v.Metadata() != "build.7" {
		t.Errorf("metadata = %q, want build.7", v.Metadata())
	}
	if v.Prefix() != "v" {
		t.Errorf("prefix = %q, want v", v.Prefix())
	}
}

func TestParseInvalid(t *testing.T) {
	cases := []string{
		"",
		"1",
		"1.2",
		"1.2.3.4",
		"01.2.3",
		"1.02.3",
		"1.2.x",
		"v",
		"1.2.3-",
		"1.2.3-01",
		"1.2.3+",
		"1.2.-3",
		"a.b.c",
	}
	for _, in := range cases {
		if _, err := version.Parse(in); err == nil {
			t.Errorf("Parse(%q) = nil error, want error", in)
		}
	}
}

// TestPrecedence covers the ordering example from semver.org §11.
func TestPrecedence(t *testing.T) {
	ordered := []string{
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"1.0.0-alpha.beta",
		"1.0.0-beta",
		"1.0.0-beta.2",
		"1.0.0-beta.11",
		"1.0.0-rc.1",
		"1.0.0",
	}
	for i := range ordered {
		for j := range ordered {
			a := version.MustParse(ordered[i])
			b := version.MustParse(ordered[j])
			got := version.Compare(a, b)
			want := cmpIndex(i, j)
			if got != want {
				t.Errorf("Compare(%q,%q) = %d, want %d", ordered[i], ordered[j], got, want)
			}
		}
	}
}

func TestCompareCoreAndBuildMetadata(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"1.0.0", "1.2.0", -1},
		{"2.0.0", "1.9.9", 1},
		{"1.0.0", "1.0.0", 0},
		{"1.0.0-rc.1", "1.0.0", -1},
		{"1.0.0+build.5", "1.0.0+build.9", 0}, // build metadata ignored
		{"1.0.0", "1.0.0+build", 0},
	}
	for _, c := range cases {
		got := version.Compare(version.MustParse(c.a), version.MustParse(c.b))
		if got != c.want {
			t.Errorf("Compare(%q,%q) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func cmpIndex(i, j int) int {
	switch {
	case i < j:
		return -1
	case i > j:
		return 1
	default:
		return 0
	}
}

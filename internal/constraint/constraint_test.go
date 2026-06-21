package constraint_test

import (
	"testing"

	"github.com/sgaunet/semver/internal/constraint"
	"github.com/sgaunet/semver/internal/version"
)

func TestCheck(t *testing.T) {
	cases := []struct {
		ver        string
		constraint string
		want       bool
	}{
		// comparators
		{"1.5.0", ">=1.2.0 <2.0.0", true},
		{"2.0.0", ">=1.2.0 <2.0.0", false},
		{"1.2.0", ">=1.2.0 <2.0.0", true},
		{"1.0.0", "!=1.0.0", false},
		{"1.0.1", "!=1.0.0", true},
		{"1.2.3", "=1.2.3", true},
		// caret
		{"1.5.0", "^1.2.3", true},
		{"2.0.0", "^1.2.3", false},
		{"0.2.5", "^0.2.3", true},
		{"0.3.0", "^0.2.3", false},
		{"0.0.3", "^0.0.3", true},
		{"0.0.4", "^0.0.3", false},
		// tilde
		{"1.2.9", "~1.2.3", true},
		{"1.3.0", "~1.2.3", false},
		{"1.9.0", "~1", true},
		{"2.0.0", "~1", false},
		// wildcards
		{"1.2.9", "1.2.x", true},
		{"1.3.0", "1.2.x", false},
		{"1.9.9", "1.x", true},
		{"2.0.0", "1.x", false},
		{"3.4.5", "*", true},
		// hyphen range
		{"1.4.0", "1.2.0 - 1.5.0", true},
		{"1.6.0", "1.2.0 - 1.5.0", false},
		{"1.5.0", "1.2.0 - 1.5.0", true},
		// OR
		{"2.3.0", "1.5.0 || 2.x", true},
		{"1.5.0", "1.5.0 || 2.x", true},
		{"3.0.0", "1.5.0 || 2.x", false},
		// pre-release exclusion
		{"1.5.0-rc.1", ">=1.2.0 <2.0.0", false},
		{"1.5.0-rc.1", ">=1.5.0-rc.0 <2.0.0", true},
	}
	for _, c := range cases {
		con, err := constraint.Parse(c.constraint)
		if err != nil {
			t.Errorf("Parse(%q): %v", c.constraint, err)
			continue
		}
		v := version.MustParse(c.ver)
		if got := con.Check(v); got != c.want {
			t.Errorf("%q satisfies %q = %v, want %v", c.ver, c.constraint, got, c.want)
		}
	}
}

func TestParseMalformed(t *testing.T) {
	for _, s := range []string{"", ">=bogus", "1.2.3.4", ">=1.2.x.y"} {
		if _, err := constraint.Parse(s); err == nil {
			t.Errorf("Parse(%q) = nil error, want error", s)
		}
	}
}

package version_test

import (
	"sort"
	"testing"

	"github.com/sgaunet/semver/internal/version"
)

// TestSortTotalOrderConsistency verifies that sorting a set agrees with pairwise
// Compare for every pair (spec SC-006).
func TestSortTotalOrderConsistency(t *testing.T) {
	raw := []string{
		"2.0.0", "1.0.0", "1.0.0-rc.1", "1.0.0-alpha", "1.2.0",
		"1.0.0-beta.2", "1.0.0-beta.11", "0.9.9", "1.0.0+build", "10.0.0",
	}
	vs := make([]version.Version, len(raw))
	for i, s := range raw {
		vs[i] = version.MustParse(s)
	}
	sort.SliceStable(vs, func(i, j int) bool {
		return version.Compare(vs[i], vs[j]) < 0
	})
	// After sorting ascending, no earlier element may be greater than a later one.
	for i := range vs {
		for j := i + 1; j < len(vs); j++ {
			if version.Compare(vs[i], vs[j]) > 0 {
				t.Errorf("order violation: %s before %s", vs[i], vs[j])
			}
		}
	}
}

// TestCompareAntisymmetry verifies Compare(a,b) == -Compare(b,a).
func TestCompareAntisymmetry(t *testing.T) {
	raw := []string{"1.0.0", "1.0.0-rc.1", "1.2.3", "2.0.0", "1.0.0-alpha.1"}
	for _, a := range raw {
		for _, b := range raw {
			va, vb := version.MustParse(a), version.MustParse(b)
			if got := version.Compare(va, vb); got != -version.Compare(vb, va) {
				t.Errorf("antisymmetry: Compare(%s,%s)=%d but Compare(%s,%s)=%d", a, b, got, b, a, version.Compare(vb, va))
			}
		}
	}
}

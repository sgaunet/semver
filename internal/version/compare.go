package version

import "strings"

// Compare returns -1, 0, or +1 as a is less than, equal to, or greater than b in
// Semantic Versioning precedence order (semver §11). Build metadata is ignored.
func Compare(a, b Version) int {
	if c := cmpUint(a.Major, b.Major); c != 0 {
		return c
	}
	if c := cmpUint(a.Minor, b.Minor); c != 0 {
		return c
	}
	if c := cmpUint(a.Patch, b.Patch); c != 0 {
		return c
	}
	return comparePre(a.Pre, b.Pre)
}

// Equal reports whether v and o have the same precedence.
func (v Version) Equal(o Version) bool { return Compare(v, o) == 0 }

// LessThan reports whether v has lower precedence than o.
func (v Version) LessThan(o Version) bool { return Compare(v, o) < 0 }

// GreaterThan reports whether v has higher precedence than o.
func (v Version) GreaterThan(o Version) bool { return Compare(v, o) > 0 }

func cmpUint(a, b uint64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

func cmpInt(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// comparePre compares pre-release identifier lists. A version with no pre-release
// has higher precedence than one that has one (when the core triple is equal).
func comparePre(a, b []string) int {
	switch {
	case len(a) == 0 && len(b) == 0:
		return 0
	case len(a) == 0:
		return 1
	case len(b) == 0:
		return -1
	}
	n := min(len(a), len(b))
	for i := range n {
		if c := comparePreIdent(a[i], b[i]); c != 0 {
			return c
		}
	}
	// All shared identifiers equal: the longer set has higher precedence.
	return cmpInt(len(a), len(b))
}

func comparePreIdent(a, b string) int {
	an, bn := isNumericIdent(a), isNumericIdent(b)
	switch {
	case an && bn:
		// No leading zeros are allowed, so a longer string is the larger number.
		if len(a) != len(b) {
			return cmpInt(len(a), len(b))
		}
		return strings.Compare(a, b)
	case an && !bn:
		// Numeric identifiers have lower precedence than alphanumeric ones.
		return -1
	case !an && bn:
		return 1
	default:
		return strings.Compare(a, b)
	}
}

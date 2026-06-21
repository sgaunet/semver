package constraint

import "github.com/sgaunet/semver/internal/version"

// match reports whether v satisfies this single comparator (precedence only; the
// pre-release-exclusion rule is applied at the group level in groupAllows).
func (c comparator) match(v version.Version) bool {
	cmp := version.Compare(v, c.v)
	switch c.o {
	case opEQ:
		return cmp == 0
	case opNE:
		return cmp != 0
	case opGT:
		return cmp > 0
	case opGE:
		return cmp >= 0
	case opLT:
		return cmp < 0
	case opLE:
		return cmp <= 0
	default:
		return false
	}
}

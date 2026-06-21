// Package constraint implements parsing and evaluation of version range
// constraints compatible with the widely used npm node-semver / Go
// Masterminds-semver syntax (comparators, caret, tilde, wildcards, hyphen ranges,
// AND/OR), including the standard pre-release-exclusion rule.
//
// It imports no CLI framework and operates on version.Version values.
package constraint

import "github.com/sgaunet/semver/internal/version"

type op int

const (
	opEQ op = iota
	opNE
	opGT
	opGE
	opLT
	opLE
)

type comparator struct {
	o op
	v version.Version
}

// Constraint is a parsed range expression: an OR of AND-groups of comparators. A
// version satisfies it when it satisfies every comparator in at least one group.
type Constraint struct {
	orGroups [][]comparator
}

// Check reports whether v satisfies the constraint.
func (c Constraint) Check(v version.Version) bool {
	for _, group := range c.orGroups {
		if groupAllows(group, v) {
			return true
		}
	}
	return false
}

// groupAllows applies the AND-group, with the pre-release-exclusion rule: a
// pre-release version satisfies a group only if some comparator in that group pins
// the same major.minor.patch and is itself a pre-release.
func groupAllows(group []comparator, v version.Version) bool {
	if len(group) == 0 {
		return false
	}
	if v.IsPrerelease() {
		allowed := false
		for _, cm := range group {
			if cm.v.IsPrerelease() && sameCore(v, cm.v) {
				allowed = true
				break
			}
		}
		if !allowed {
			return false
		}
	}
	for _, cm := range group {
		if !cm.match(v) {
			return false
		}
	}
	return true
}

func sameCore(a, b version.Version) bool {
	return a.Major == b.Major && a.Minor == b.Minor && a.Patch == b.Patch
}

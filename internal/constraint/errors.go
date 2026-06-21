package constraint

import "errors"

// Sentinel errors returned (sometimes wrapped with positional context) while parsing
// constraint expressions. Defined here to satisfy the err113 "no dynamic errors" rule
// and to let callers match them with errors.Is.
var (
	errEmptyConstraint  = errors.New("empty constraint")
	errEmptyGroup       = errors.New("empty constraint group")
	errUnexpectedHyphen = errors.New("unexpected '-' in constraint")
	errNoComparators    = errors.New("no comparators in constraint group")
	errEmptyVersion     = errors.New("empty version in constraint")
	errInvalidVersion   = errors.New("invalid version")
	errInvalidComponent = errors.New("invalid version component")
)

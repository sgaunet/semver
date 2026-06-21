// Package version implements parsing, comparison, and manipulation of semantic
// versions per the Semantic Versioning 2.0.0 specification (https://semver.org).
//
// It imports no CLI framework; all I/O, formatting, and exit-code concerns live in
// the cli package. Operations never mutate their receiver — every transformation
// returns a new Version value.
package version

import "strings"

// Version is a parsed semantic version. The zero value is not a meaningful version;
// obtain values via Parse or MustParse.
type Version struct {
	Major uint64
	Minor uint64
	Patch uint64
	// Pre holds the ordered pre-release identifiers (the part after '-'); an empty
	// slice means a stable release.
	Pre []string
	// Build holds the ordered build-metadata identifiers (the part after '+'); it is
	// ignored when determining precedence (semver §10).
	Build []string
	// prefix is the original leading prefix as written: "", "v", or "V". It is a
	// presentation detail re-emitted by String and carried across operations.
	prefix string
}

// IsPrerelease reports whether v carries a pre-release identifier.
func (v Version) IsPrerelease() bool { return len(v.Pre) > 0 }

// Prerelease returns the dot-joined pre-release identifiers, or "" if none.
func (v Version) Prerelease() string { return strings.Join(v.Pre, ".") }

// Metadata returns the dot-joined build-metadata identifiers, or "" if none.
func (v Version) Metadata() string { return strings.Join(v.Build, ".") }

// Prefix returns the original leading prefix ("", "v", or "V").
func (v Version) Prefix() string { return v.prefix }

// clone returns a deep copy so that mutations of slice fields do not alias.
func (v Version) clone() Version {
	out := v
	if v.Pre != nil {
		out.Pre = append([]string(nil), v.Pre...)
	}
	if v.Build != nil {
		out.Build = append([]string(nil), v.Build...)
	}
	return out
}

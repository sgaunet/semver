package version

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Parse parses s as a semantic version, accepting an optional leading 'v' or 'V'
// prefix which is preserved by String. It applies the strict Semantic Versioning
// 2.0.0 grammar and returns an error naming the first problem found, suitable for
// printing to stderr.
func Parse(s string) (Version, error) {
	if s == "" {
		return Version{}, errors.New(`invalid version "": empty string`)
	}
	orig := s
	var v Version
	if s[0] == 'v' || s[0] == 'V' {
		v.prefix = s[:1]
		s = s[1:]
	}

	// Build metadata: everything after the first '+'.
	core := s
	if before, after, found := strings.Cut(s, "+"); found {
		build, err := parseIdentifiers(after, true)
		if err != nil {
			return Version{}, fmt.Errorf("invalid version %q: build metadata: %w", orig, err)
		}
		v.Build = build
		core = before
	}

	// Pre-release: everything after the first '-' in the core.
	numbers := core
	if before, after, found := strings.Cut(core, "-"); found {
		pre, err := parseIdentifiers(after, false)
		if err != nil {
			return Version{}, fmt.Errorf("invalid version %q: pre-release: %w", orig, err)
		}
		v.Pre = pre
		numbers = before
	}

	// Core: major.minor.patch.
	parts := strings.Split(numbers, ".")
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("invalid version %q: expected major.minor.patch", orig)
	}
	names := [3]string{"major", "minor", "patch"}
	var out [3]uint64
	for i, p := range parts {
		n, err := parseNumeric(p)
		if err != nil {
			return Version{}, fmt.Errorf("invalid version %q: %s %w", orig, names[i], err)
		}
		out[i] = n
	}
	v.Major, v.Minor, v.Patch = out[0], out[1], out[2]
	return v, nil
}

// MustParse is like Parse but panics on error. Intended for tests and constants.
func MustParse(s string) Version {
	v, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return v
}

func parseNumeric(s string) (uint64, error) {
	if s == "" {
		return 0, errors.New("component is missing")
	}
	if len(s) > 1 && s[0] == '0' {
		return 0, fmt.Errorf("component has a leading zero (%q)", s)
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0, fmt.Errorf("component is not numeric (%q)", s)
		}
	}
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("component %q is out of range", s)
	}
	return n, nil
}

// parseIdentifiers validates a dot-separated identifier list. For pre-release
// identifiers (build == false), a purely numeric identifier must not have a leading
// zero; build-metadata identifiers may.
func parseIdentifiers(s string, build bool) ([]string, error) {
	if s == "" {
		return nil, errors.New("is empty")
	}
	ids := strings.Split(s, ".")
	for _, id := range ids {
		if id == "" {
			return nil, errors.New("has an empty identifier")
		}
		if !isAlnumHyphen(id) {
			return nil, fmt.Errorf("has invalid identifier %q", id)
		}
		if !build && isNumericIdent(id) && len(id) > 1 && id[0] == '0' {
			return nil, fmt.Errorf("has a leading zero in numeric identifier %q", id)
		}
	}
	return ids, nil
}

func isAlnumHyphen(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= '0' && c <= '9', c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c == '-':
		default:
			return false
		}
	}
	return true
}

func isNumericIdent(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

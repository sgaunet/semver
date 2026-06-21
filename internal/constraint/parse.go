package constraint

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sgaunet/semver/internal/version"
)

// Parse parses a constraint expression. Groups separated by "||" are ORed; within a
// group, comparators separated by spaces or commas are ANDed. Supported atoms:
// comparators (=, !=, >, >=, <, <=), caret (^), tilde (~), wildcards (x/X/*),
// and hyphen ranges (A - B).
func Parse(s string) (Constraint, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Constraint{}, errEmptyConstraint
	}
	var c Constraint
	for orPart := range strings.SplitSeq(s, "||") {
		group, err := parseGroup(strings.TrimSpace(orPart))
		if err != nil {
			return Constraint{}, err
		}
		c.orGroups = append(c.orGroups, group)
	}
	return c, nil
}

func parseGroup(s string) ([]comparator, error) {
	if s == "" {
		return nil, errEmptyGroup
	}
	tokens := strings.FieldsFunc(s, isSeparator)
	group, err := comparatorsFromTokens(tokens)
	if err != nil {
		return nil, err
	}
	if len(group) == 0 {
		return nil, errNoComparators
	}
	return group, nil
}

// comparatorsFromTokens turns AND-group tokens into comparators, expanding hyphen
// ranges ("A - B") spread across three tokens.
func comparatorsFromTokens(tokens []string) ([]comparator, error) {
	var group []comparator
	for i := 0; i < len(tokens); i++ {
		// Hyphen range: "A - B".
		if i+2 < len(tokens) && tokens[i+1] == "-" {
			cmps, err := hyphenRange(tokens[i], tokens[i+2])
			if err != nil {
				return nil, err
			}
			group = append(group, cmps...)
			i += 2
			continue
		}
		if tokens[i] == "-" {
			return nil, errUnexpectedHyphen
		}
		cmps, err := parseComparator(tokens[i])
		if err != nil {
			return nil, err
		}
		group = append(group, cmps...)
	}
	return group, nil
}

// isSeparator reports whether r delimits comparators within an AND-group.
func isSeparator(r rune) bool {
	return r == ' ' || r == '\t' || r == ','
}

func parseComparator(tok string) ([]comparator, error) {
	switch {
	case strings.HasPrefix(tok, "^"):
		return caret(tok[1:])
	case strings.HasPrefix(tok, "~"):
		return tilde(tok[1:])
	}

	o, rest := splitOperator(tok)
	p, err := parsePartial(rest)
	if err != nil {
		return nil, err
	}
	if o == opEQ && (p.wildcard || !p.full()) {
		// A bare partial or wildcard equality expands to a range.
		return p.rangeComparators(), nil
	}
	return []comparator{{o: o, v: p.version()}}, nil
}

// splitOperator strips a leading comparison operator from tok, returning the operator
// (defaulting to equality) and the remaining version operand.
func splitOperator(tok string) (op, string) {
	switch {
	case strings.HasPrefix(tok, ">="):
		return opGE, tok[2:]
	case strings.HasPrefix(tok, "<="):
		return opLE, tok[2:]
	case strings.HasPrefix(tok, "!="):
		return opNE, tok[2:]
	case strings.HasPrefix(tok, ">"):
		return opGT, tok[1:]
	case strings.HasPrefix(tok, "<"):
		return opLT, tok[1:]
	case strings.HasPrefix(tok, "=="):
		return opEQ, tok[2:]
	case strings.HasPrefix(tok, "="):
		return opEQ, tok[1:]
	default:
		return opEQ, tok
	}
}

func caret(s string) ([]comparator, error) {
	p, err := parsePartial(s)
	if err != nil {
		return nil, err
	}
	var hi version.Version
	switch {
	case p.major != 0:
		hi = ver(p.major+1, 0, 0)
	case p.minor != 0:
		hi = ver(0, p.minor+1, 0)
	case p.patch != 0:
		hi = ver(0, 0, p.patch+1)
	case p.hasPatch:
		hi = ver(0, 0, 1) // ^0.0.0
	case p.hasMinor:
		hi = ver(0, 1, 0) // ^0.0
	default:
		hi = ver(1, 0, 0) // ^0 or ^0.x
	}
	return []comparator{{opGE, p.version()}, {opLT, hi}}, nil
}

func tilde(s string) ([]comparator, error) {
	p, err := parsePartial(s)
	if err != nil {
		return nil, err
	}
	var hi version.Version
	if p.hasMinor {
		hi = ver(p.major, p.minor+1, 0)
	} else {
		hi = ver(p.major+1, 0, 0)
	}
	return []comparator{{opGE, p.version()}, {opLT, hi}}, nil
}

func hyphenRange(loTok, hiTok string) ([]comparator, error) {
	lo, err := parsePartial(loTok)
	if err != nil {
		return nil, err
	}
	hi, err := parsePartial(hiTok)
	if err != nil {
		return nil, err
	}
	cmps := []comparator{{opGE, ver(lo.major, lo.minor, lo.patch)}}
	switch {
	case hi.hasPatch:
		cmps = append(cmps, comparator{opLE, ver(hi.major, hi.minor, hi.patch)})
	case hi.hasMinor:
		cmps = append(cmps, comparator{opLT, ver(hi.major, hi.minor+1, 0)})
	default:
		cmps = append(cmps, comparator{opLT, ver(hi.major+1, 0, 0)})
	}
	return cmps, nil
}

// partial is a possibly-incomplete version operand from a constraint.
type partial struct {
	major, minor, patch          uint64
	hasMajor, hasMinor, hasPatch bool
	wildcard                     bool
	pre, build                   []string
}

func (p partial) full() bool {
	return p.hasMajor && p.hasMinor && p.hasPatch && !p.wildcard
}

func (p partial) version() version.Version {
	return version.Version{Major: p.major, Minor: p.minor, Patch: p.patch, Pre: p.pre, Build: p.build}
}

// rangeComparators expands a bare partial/wildcard into bound comparators.
func (p partial) rangeComparators() []comparator {
	switch {
	case !p.hasMajor: // "*" or "x" — match anything (>= 0.0.0)
		return []comparator{{opGE, ver(0, 0, 0)}}
	case !p.hasMinor: // "1", "1.x" — >=1.0.0 <2.0.0
		return []comparator{{opGE, ver(p.major, 0, 0)}, {opLT, ver(p.major+1, 0, 0)}}
	case !p.hasPatch: // "1.2", "1.2.x" — >=1.2.0 <1.3.0
		return []comparator{{opGE, ver(p.major, p.minor, 0)}, {opLT, ver(p.major, p.minor+1, 0)}}
	default:
		return []comparator{{opEQ, p.version()}}
	}
}

// Indices of the dotted segments of a version operand, and the maximum number of them.
const (
	majorIndex  = 0
	minorIndex  = 1
	patchIndex  = 2
	maxSegments = 3
)

func parsePartial(s string) (partial, error) {
	var p partial
	if s == "" {
		return p, errEmptyVersion
	}
	if s[0] == 'v' || s[0] == 'V' {
		s = s[1:]
	}
	if before, after, found := strings.Cut(s, "+"); found {
		p.build = splitDot(after)
		s = before
	}
	if before, after, found := strings.Cut(s, "-"); found {
		p.pre = splitDot(after)
		s = before
	}
	segs := strings.Split(s, ".")
	if len(segs) > maxSegments {
		return p, fmt.Errorf("%w %q in constraint", errInvalidVersion, s)
	}
	if err := parseSegments(&p, segs); err != nil {
		return p, err
	}
	return p, nil
}

// parseSegments fills p's numeric components from the dotted segments, stopping at the
// first wildcard ('x', 'X', or '*'), after which the remaining components are
// unspecified.
func parseSegments(p *partial, segs []string) error {
	for idx, seg := range segs {
		if seg == "x" || seg == "X" || seg == "*" {
			p.wildcard = true
			break
		}
		n, err := strconv.ParseUint(seg, 10, 64)
		if err != nil {
			return fmt.Errorf("%w %q in constraint", errInvalidComponent, seg)
		}
		switch idx {
		case majorIndex:
			p.major, p.hasMajor = n, true
		case minorIndex:
			p.minor, p.hasMinor = n, true
		case patchIndex:
			p.patch, p.hasPatch = n, true
		}
	}
	return nil
}

func ver(major, minor, patch uint64) version.Version {
	return version.Version{Major: major, Minor: minor, Patch: patch}
}

func splitDot(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ".")
}

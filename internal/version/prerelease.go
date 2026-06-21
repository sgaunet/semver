package version

import (
	"fmt"
	"strconv"
)

// StartPrerelease applies the given core bump and begins a pre-release series with
// the supplied identifier — e.g. v1.0.0.StartPrerelease("rc", BumpMinor) yields
// v1.1.0-rc.1. The identifier must be a valid pre-release identifier. Build metadata
// is cleared.
func (v Version) StartPrerelease(id string, k BumpKind) (Version, error) {
	if id == "" {
		return Version{}, errPreIdentEmpty
	}
	if !isValidPreIdentifier(id) {
		return Version{}, fmt.Errorf("%w %q", errPreIdentInvalid, id)
	}
	out := v.bump(k)
	out.Pre = []string{id, "1"}
	out.Build = nil
	return out, nil
}

// IncPrerelease increments the numeric tail of the pre-release (rc.1 -> rc.2), or
// appends ".1" when the last identifier is not numeric (beta -> beta.1). It returns
// an error if v has no pre-release. Build metadata is cleared.
func (v Version) IncPrerelease() (Version, error) {
	if !v.IsPrerelease() {
		return Version{}, errNoPreToIncrement
	}
	out := v.clone()
	out.Build = nil
	last := out.Pre[len(out.Pre)-1]
	if isNumericIdent(last) {
		n, err := strconv.ParseUint(last, 10, 64)
		if err != nil {
			return Version{}, fmt.Errorf("%w (%q)", errPreCounterRange, last)
		}
		out.Pre[len(out.Pre)-1] = strconv.FormatUint(n+1, 10)
	} else {
		out.Pre = append(out.Pre, "1")
	}
	return out, nil
}

// Finalize promotes a pre-release to its stable release by dropping the pre-release
// identifiers (e.g. v1.1.0-rc.1 -> v1.1.0). It returns an error if v has no
// pre-release. Build metadata is cleared.
func (v Version) Finalize() (Version, error) {
	if !v.IsPrerelease() {
		return Version{}, errNoPreToFinalize
	}
	out := v.clone()
	out.Pre = nil
	out.Build = nil
	return out, nil
}

func isValidPreIdentifier(id string) bool {
	if !isAlnumHyphen(id) {
		return false
	}
	if isNumericIdent(id) && len(id) > 1 && id[0] == '0' {
		return false
	}
	return true
}

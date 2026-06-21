package version

import (
	"strconv"
	"strings"
)

// String renders the canonical semver text, preserving the original 'v'/'V' prefix
// the value was parsed with. The result round-trips: Parse(v.String()) equals v.
func (v Version) String() string {
	var b strings.Builder
	b.WriteString(v.prefix)
	b.WriteString(strconv.FormatUint(v.Major, 10))
	b.WriteByte('.')
	b.WriteString(strconv.FormatUint(v.Minor, 10))
	b.WriteByte('.')
	b.WriteString(strconv.FormatUint(v.Patch, 10))
	if len(v.Pre) > 0 {
		b.WriteByte('-')
		b.WriteString(strings.Join(v.Pre, "."))
	}
	if len(v.Build) > 0 {
		b.WriteByte('+')
		b.WriteString(strings.Join(v.Build, "."))
	}
	return b.String()
}

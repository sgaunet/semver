package cli

import "errors"

// Output formats accepted by --output.
const (
	formatText = "text"
	formatJSON = "json"
)

// Core bump component names, used both as command names and as --bump/component values.
const (
	cmdMajor = "major"
	cmdMinor = "minor"
	cmdPatch = "patch"
)

// twoArgs is the positional-argument count for the two-operand commands (compare, get,
// satisfies).
const twoArgs = 2

// Scanner buffer sizes for reading versions from stdin: a 64 KiB initial buffer that
// grows up to 1 MiB per line.
const (
	scanBufSize    = 64 * 1024
	scanBufMaxSize = 1024 * 1024
)

// errInvalidBump is returned when --bump names something other than a core component.
var errInvalidBump = errors.New("want major, minor, or patch")

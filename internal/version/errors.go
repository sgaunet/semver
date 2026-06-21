package version

import "errors"

// Sentinel errors returned (often wrapped with positional context) by parsing and
// pre-release operations. They are defined here so callers can match them with
// errors.Is and to satisfy the err113 "no dynamic errors" rule.
var (
	errEmptyVersion         = errors.New("empty string")
	errFormat               = errors.New("expected major.minor.patch")
	errComponentMissing     = errors.New("component is missing")
	errComponentLeadingZero = errors.New("component has a leading zero")
	errComponentNotNumeric  = errors.New("component is not numeric")
	errComponentRange       = errors.New("component is out of range")
	errIdentListEmpty       = errors.New("is empty")
	errIdentEmpty           = errors.New("has an empty identifier")
	errIdentInvalid         = errors.New("has invalid identifier")
	errIdentLeadingZero     = errors.New("has a leading zero in numeric identifier")

	errPreIdentEmpty    = errors.New("pre-release identifier is empty")
	errPreIdentInvalid  = errors.New("invalid pre-release identifier")
	errNoPreToIncrement = errors.New("version has no pre-release to increment")
	errPreCounterRange  = errors.New("pre-release counter is out of range")
	errNoPreToFinalize  = errors.New("version has no pre-release to finalize")
)

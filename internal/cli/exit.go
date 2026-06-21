package cli

// Exit codes. See specs/001-semver-cli/contracts/cli-contract.md for the
// authoritative registry. The codes compose in shell conditionals.
const (
	// CodeOK is success / an affirmative result (equal, valid, satisfied, present).
	CodeOK = 0
	// CodeError is a generic or internal failure (unexpected I/O, encoding, etc.).
	CodeError = 1
	// CodeUsage is a usage error: bad flags or arguments, a malformed constraint, or
	// an invalid version supplied to a command that requires a valid one.
	CodeUsage = 2
	// CodeLower is the "less than" comparison result, and the shared negative-domain
	// code reused by the boolean queries below.
	CodeLower = 10
	// CodeHigher is the "greater than" comparison result (compare only).
	CodeHigher = 11
)

// Negative-domain results reuse CodeLower (10).
const (
	CodeNotSatisfied = CodeLower // satisfies: false
	CodeInvalid      = CodeLower // validate: invalid
	CodeAbsent       = CodeLower // get: component absent
)

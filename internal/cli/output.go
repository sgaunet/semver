package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// streams bundles a command's I/O endpoints and presentation options. stdout
// carries machine-parseable data only; stderr carries diagnostics.
type streams struct {
	in      io.Reader
	out     io.Writer
	err     io.Writer
	format  string // "text" or "json"
	quiet   bool
	verbose bool
	noColor bool
}

// emit writes the command result: the bare token in text mode, or the JSON encoding
// of jsonVal in json mode. It returns CodeOK, or CodeError if encoding fails.
func (s streams) emit(text string, jsonVal any) int {
	if s.format == formatJSON {
		b, err := json.Marshal(jsonVal)
		if err != nil {
			s.errorf("encoding json: %v", err)
			return CodeError
		}
		fmt.Fprintln(s.out, string(b))
		return CodeOK
	}
	fmt.Fprintln(s.out, text)
	return CodeOK
}

// errorf writes a diagnostic to stderr, prefixed with the program name. Errors are
// always shown; --quiet suppresses only non-error output.
func (s streams) errorf(format string, args ...any) {
	fmt.Fprintf(s.err, "semver: "+format+"\n", args...)
}

// verbosef writes detail to stderr only when --verbose is set and --quiet is not.
func (s streams) verbosef(format string, args ...any) {
	if s.verbose && !s.quiet {
		fmt.Fprintf(s.err, format+"\n", args...)
	}
}

func noColorFromEnv() bool {
	return os.Getenv("NO_COLOR") != ""
}

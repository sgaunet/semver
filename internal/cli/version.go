package cli

import (
	"context"
	"fmt"
)

type versionInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

func runVersion(_ context.Context, args []string, s streams) int {
	// version accepts the global flags (notably --output) but no positional args.
	fs := setup("version", "", "Print build version information", &s)
	rest, code, ok := parse(fs, args, &s)
	if !ok {
		return code
	}
	if len(rest) != 0 {
		s.errorf("version takes no arguments")
		return CodeUsage
	}
	text := fmt.Sprintf("semver %s (%s, %s)", Version, Commit, Date)
	return s.emit(text, versionInfo{Version: Version, Commit: Commit, Date: Date})
}

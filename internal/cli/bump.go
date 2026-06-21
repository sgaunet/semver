package cli

import (
	"context"

	"github.com/sgaunet/semver/internal/version"
)

type bumpResult struct {
	Input     string `json:"input"`
	Operation string `json:"operation"`
	Result    string `json:"result"`
}

// runBump returns the handler for the major/minor/patch commands.
func runBump(op string) commandFunc {
	return func(_ context.Context, args []string, s streams) int {
		fs := setup(op, "<version>", "Bump the "+op+" version", &s)
		rest, code, ok := parse(fs, args, &s)
		if !ok {
			return code
		}
		if len(rest) != 1 {
			s.errorf("expected exactly one version argument")
			return CodeUsage
		}
		v, err := version.Parse(rest[0])
		if err != nil {
			s.errorf("%v", err)
			return CodeUsage
		}
		var res version.Version
		switch op {
		case cmdMajor:
			res = v.IncMajor()
		case cmdMinor:
			res = v.IncMinor()
		default:
			res = v.IncPatch()
		}
		return s.emit(res.String(), bumpResult{Input: rest[0], Operation: op, Result: res.String()})
	}
}

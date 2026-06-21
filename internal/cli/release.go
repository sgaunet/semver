package cli

import (
	"context"

	"github.com/sgaunet/semver/internal/version"
)

type releaseResult struct {
	Input  string `json:"input"`
	Result string `json:"result"`
}

func runRelease(_ context.Context, args []string, s streams) int {
	fs := setup("release", "<version>", "Finalize a pre-release to its stable version", &s)
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
	res, err := v.Finalize()
	if err != nil {
		s.errorf("%v", err)
		return CodeUsage
	}
	return s.emit(res.String(), releaseResult{Input: rest[0], Result: res.String()})
}

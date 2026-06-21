package cli

import (
	"context"
	"fmt"

	"github.com/sgaunet/semver/internal/version"
)

type prereleaseResult struct {
	Input      string `json:"input"`
	Result     string `json:"result"`
	Prerelease string `json:"prerelease,omitempty"`
}

func runPrerelease(_ context.Context, args []string, s streams) int {
	var pre, bump string
	fs := setup("prerelease", "<version>", "Start or increment a pre-release", &s)
	fs.StringVar(&pre, "pre", "", "pre-release identifier to start (e.g. rc, beta)")
	fs.StringVar(&bump, "bump", "patch", "core bump when starting: major|minor|patch")
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
	if pre != "" {
		kind, kerr := bumpKind(bump)
		if kerr != nil {
			s.errorf("%v", kerr)
			return CodeUsage
		}
		res, err = v.StartPrerelease(pre, kind)
	} else {
		res, err = v.IncPrerelease()
	}
	if err != nil {
		s.errorf("%v", err)
		return CodeUsage
	}
	return s.emit(res.String(), prereleaseResult{
		Input:      rest[0],
		Result:     res.String(),
		Prerelease: res.Prerelease(),
	})
}

func bumpKind(s string) (version.BumpKind, error) {
	switch s {
	case "major":
		return version.BumpMajor, nil
	case "minor":
		return version.BumpMinor, nil
	case "patch":
		return version.BumpPatch, nil
	default:
		return 0, fmt.Errorf("invalid --bump %q: want major, minor, or patch", s)
	}
}

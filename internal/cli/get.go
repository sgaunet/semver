package cli

import (
	"context"
	"strconv"

	"github.com/sgaunet/semver/internal/version"
)

type getResult struct {
	Version   string `json:"version"`
	Component string `json:"component"`
	Value     string `json:"value"`
	Present   bool   `json:"present"`
}

func runGet(_ context.Context, args []string, s streams) int {
	fs := setup("get", "<component> <version>", "Extract a component of a version", &s)
	rest, code, ok := parse(fs, args, &s)
	if !ok {
		return code
	}
	if len(rest) != 2 {
		s.errorf("expected: get <component> <version>")
		return CodeUsage
	}
	comp := rest[0]
	v, err := version.Parse(rest[1])
	if err != nil {
		s.errorf("%v", err)
		return CodeUsage
	}

	var value string
	present := true
	switch comp {
	case "major":
		value = strconv.FormatUint(v.Major, 10)
	case "minor":
		value = strconv.FormatUint(v.Minor, 10)
	case "patch":
		value = strconv.FormatUint(v.Patch, 10)
	case "prerelease":
		value, present = v.Prerelease(), v.IsPrerelease()
	case "build":
		value, present = v.Metadata(), len(v.Build) > 0
	default:
		s.errorf("unknown component %q: want major, minor, patch, prerelease, or build", comp)
		return CodeUsage
	}

	if rc := s.emit(value, getResult{Version: rest[1], Component: comp, Value: value, Present: present}); rc != CodeOK {
		return rc
	}
	if !present {
		return CodeAbsent
	}
	return CodeOK
}

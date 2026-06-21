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
	if len(rest) != twoArgs {
		s.errorf("expected: get <component> <version>")
		return CodeUsage
	}
	comp := rest[0]
	v, err := version.Parse(rest[1])
	if err != nil {
		s.errorf("%v", err)
		return CodeUsage
	}

	value, present, known := getComponent(comp, v)
	if !known {
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

// getComponent extracts the named component from v. It returns the component's string
// value, whether it is present on v, and whether comp names a known component.
func getComponent(comp string, v version.Version) (string, bool, bool) {
	switch comp {
	case cmdMajor:
		return strconv.FormatUint(v.Major, 10), true, true
	case cmdMinor:
		return strconv.FormatUint(v.Minor, 10), true, true
	case cmdPatch:
		return strconv.FormatUint(v.Patch, 10), true, true
	case "prerelease":
		return v.Prerelease(), v.IsPrerelease(), true
	case "build":
		return v.Metadata(), len(v.Build) > 0, true
	default:
		return "", false, false
	}
}

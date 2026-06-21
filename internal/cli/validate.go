package cli

import (
	"context"

	"github.com/sgaunet/semver/internal/version"
)

type validateOK struct {
	Input      string `json:"input"`
	Valid      bool   `json:"valid"`
	Major      uint64 `json:"major"`
	Minor      uint64 `json:"minor"`
	Patch      uint64 `json:"patch"`
	Prerelease string `json:"prerelease"`
	Build      string `json:"build"`
}

type validateErr struct {
	Input string `json:"input"`
	Valid bool   `json:"valid"`
	Error string `json:"error"`
}

func runValidate(_ context.Context, args []string, s streams) int {
	fs := setup("validate", "<version>", "Validate and parse a version", &s)
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
		if s.format == formatJSON {
			s.emit("", validateErr{Input: rest[0], Valid: false, Error: err.Error()})
		} else {
			s.errorf("%v", err)
		}
		return CodeInvalid
	}

	if s.format == formatJSON {
		return s.emit("", validateOK{
			Input:      rest[0],
			Valid:      true,
			Major:      v.Major,
			Minor:      v.Minor,
			Patch:      v.Patch,
			Prerelease: v.Prerelease(),
			Build:      v.Metadata(),
		})
	}
	return s.emit("valid", nil)
}

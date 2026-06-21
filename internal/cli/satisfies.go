package cli

import (
	"context"

	"github.com/sgaunet/semver/internal/constraint"
	"github.com/sgaunet/semver/internal/version"
)

type satisfiesResult struct {
	Version    string `json:"version"`
	Constraint string `json:"constraint"`
	Satisfied  bool   `json:"satisfied"`
}

func runSatisfies(_ context.Context, args []string, s streams) int {
	fs := setup("satisfies", "<version> <constraint>", "Test a version against a constraint", &s)
	rest, code, ok := parse(fs, args, &s)
	if !ok {
		return code
	}
	if len(rest) != twoArgs {
		s.errorf("expected: satisfies <version> <constraint>")
		return CodeUsage
	}
	v, err := version.Parse(rest[0])
	if err != nil {
		s.errorf("%v", err)
		return CodeUsage
	}
	c, err := constraint.Parse(rest[1])
	if err != nil {
		s.errorf("%v", err)
		return CodeUsage
	}

	satisfied := c.Check(v)
	word := "false"
	if satisfied {
		word = "true"
	}
	if rc := s.emit(word, satisfiesResult{Version: rest[0], Constraint: rest[1], Satisfied: satisfied}); rc != CodeOK {
		return rc
	}
	if satisfied {
		return CodeOK
	}
	return CodeNotSatisfied
}

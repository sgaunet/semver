package cli

import (
	"context"

	"github.com/sgaunet/semver/internal/version"
)

type compareResult struct {
	A      string `json:"a"`
	B      string `json:"b"`
	Result string `json:"result"`
	Cmp    int    `json:"cmp"`
}

func runCompare(_ context.Context, args []string, s streams) int {
	fs := setup("compare", "<a> <b>", "Compare two versions", &s)
	rest, code, ok := parse(fs, args, &s)
	if !ok {
		return code
	}
	if len(rest) != twoArgs {
		s.errorf("expected exactly two version arguments")
		return CodeUsage
	}
	a, err := version.Parse(rest[0])
	if err != nil {
		s.errorf("%v", err)
		return CodeUsage
	}
	b, err := version.Parse(rest[1])
	if err != nil {
		s.errorf("%v", err)
		return CodeUsage
	}

	cmp := version.Compare(a, b)
	word := comparisonWord(cmp)
	if rc := s.emit(word, compareResult{A: rest[0], B: rest[1], Result: word, Cmp: cmp}); rc != CodeOK {
		return rc
	}
	switch {
	case cmp < 0:
		return CodeLower
	case cmp > 0:
		return CodeHigher
	default:
		return CodeOK
	}
}

func comparisonWord(cmp int) string {
	switch {
	case cmp < 0:
		return "lower"
	case cmp > 0:
		return "higher"
	default:
		return "equal"
	}
}

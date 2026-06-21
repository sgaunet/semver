package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/sgaunet/semver/internal/version"
)

type sortResult struct {
	Order    string   `json:"order"`
	Versions []string `json:"versions"`
}

// sortItem pairs an input string with its parsed version so the original spelling can
// be re-emitted after sorting by precedence.
type sortItem struct {
	raw string
	v   version.Version
}

func runSort(ctx context.Context, args []string, s streams) int {
	var desc bool
	fs := setup("sort", "[versions...]", "Sort versions by precedence", &s)
	fs.BoolVar(&desc, "desc", false, "sort in descending order")
	inputs, code, ok := parse(fs, args, &s)
	if !ok {
		return code
	}

	if len(inputs) == 0 {
		lines, rc, gotInput := readStdin(ctx, s)
		if !gotInput {
			return rc
		}
		inputs = lines
	}

	items, rc, parsed := parseVersions(inputs, s)
	if !parsed {
		return rc
	}
	s.verbosef("sorting %d version(s)", len(items))

	sortItems(items, desc)

	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.raw
	}
	return emitSorted(s, out, desc)
}

// readStdin reads versions from stdin, mapping failures to an exit code. The bool is
// false when the caller should return rc instead of proceeding.
func readStdin(ctx context.Context, s streams) ([]string, int, bool) {
	lines, err := readLines(ctx, s.in)
	if err != nil {
		if ctx.Err() != nil {
			s.errorf("cancelled")
		} else {
			s.errorf("reading stdin: %v", err)
		}
		return nil, CodeError, false
	}
	return lines, CodeOK, true
}

// parseVersions parses every input, reporting the first invalid one. The bool is false
// when the caller should return rc instead of proceeding.
func parseVersions(inputs []string, s streams) ([]sortItem, int, bool) {
	items := make([]sortItem, 0, len(inputs))
	for _, in := range inputs {
		v, err := version.Parse(in)
		if err != nil {
			s.errorf("%v", err)
			return nil, CodeUsage, false
		}
		items = append(items, sortItem{raw: in, v: v})
	}
	return items, CodeOK, true
}

// sortItems orders items by semver precedence, ascending unless desc is set.
func sortItems(items []sortItem, desc bool) {
	sort.SliceStable(items, func(i, j int) bool {
		c := version.Compare(items[i].v, items[j].v)
		if desc {
			return c > 0
		}
		return c < 0
	})
}

// emitSorted writes the sorted versions as JSON or one per line.
func emitSorted(s streams, out []string, desc bool) int {
	order := "asc"
	if desc {
		order = "desc"
	}
	if s.format == formatJSON {
		return s.emit("", sortResult{Order: order, Versions: out})
	}
	for _, line := range out {
		fmt.Fprintln(s.out, line)
	}
	return CodeOK
}

// readLines reads whitespace-trimmed, non-empty lines from r, honoring ctx so the
// read can be cancelled on SIGINT/SIGTERM.
func readLines(ctx context.Context, r io.Reader) ([]string, error) {
	type result struct {
		lines []string
		err   error
	}
	ch := make(chan result, 1)
	go func() {
		var lines []string
		sc := bufio.NewScanner(r)
		sc.Buffer(make([]byte, 0, scanBufSize), scanBufMaxSize)
		for sc.Scan() {
			if line := strings.TrimSpace(sc.Text()); line != "" {
				lines = append(lines, line)
			}
		}
		ch <- result{lines: lines, err: sc.Err()}
	}()
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("read cancelled: %w", ctx.Err())
	case res := <-ch:
		return res.lines, res.err
	}
}

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

func runSort(ctx context.Context, args []string, s streams) int {
	var desc bool
	fs := setup("sort", "[versions...]", "Sort versions by precedence", &s)
	fs.BoolVar(&desc, "desc", false, "sort in descending order")
	inputs, code, ok := parse(fs, args, &s)
	if !ok {
		return code
	}

	if len(inputs) == 0 {
		lines, err := readLines(ctx, s.in)
		if err != nil {
			if ctx.Err() != nil {
				s.errorf("cancelled")
			} else {
				s.errorf("reading stdin: %v", err)
			}
			return CodeError
		}
		inputs = lines
	}

	type item struct {
		raw string
		v   version.Version
	}
	items := make([]item, 0, len(inputs))
	for _, in := range inputs {
		v, err := version.Parse(in)
		if err != nil {
			s.errorf("%v", err)
			return CodeUsage
		}
		items = append(items, item{raw: in, v: v})
	}
	s.verbosef("sorting %d version(s)", len(items))

	sort.SliceStable(items, func(i, j int) bool {
		c := version.Compare(items[i].v, items[j].v)
		if desc {
			return c > 0
		}
		return c < 0
	})

	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.raw
	}

	order := "asc"
	if desc {
		order = "desc"
	}
	if s.format == "json" {
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
		sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for sc.Scan() {
			if line := strings.TrimSpace(sc.Text()); line != "" {
				lines = append(lines, line)
			}
		}
		ch <- result{lines: lines, err: sc.Err()}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-ch:
		return res.lines, res.err
	}
}

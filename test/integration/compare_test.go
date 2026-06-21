package integration

import "testing"

func TestIntegrationCompare(t *testing.T) {
	cases := []struct {
		a, b, word string
		code       int
	}{
		{"v1.0.0", "v1.2.0", "lower", 10},
		{"v2.0.0", "v1.9.9", "higher", 11},
		{"v1.0.0", "v1.0.0", "equal", 0},
		{"v1.0.0-rc.1", "v1.0.0", "lower", 10},           // pre-release < release
		{"v1.0.0+build.5", "v1.0.0+build.9", "equal", 0}, // build metadata ignored
	}
	for _, c := range cases {
		out, _, code := run(t, "", "compare", c.a, c.b)
		if code != c.code || out != c.word+"\n" {
			t.Errorf("compare %s %s: out=%q code=%d, want %q/%d", c.a, c.b, out, code, c.word, c.code)
		}
	}
}

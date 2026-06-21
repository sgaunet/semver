package integration

import "testing"

func TestIntegrationSatisfies(t *testing.T) {
	cases := []struct {
		ver, constraint, word string
		code                  int
	}{
		{"v1.5.0", "^1.2.0", "true", 0},
		{"v1.5.0", ">=1.2.0 <2.0.0", "true", 0},
		{"v2.0.0", ">=1.2.0 <2.0.0", "false", 10},
		{"v1.5.0-rc.1", ">=1.2.0 <2.0.0", "false", 10}, // pre-release exclusion
	}
	for _, c := range cases {
		out, _, code := run(t, "", "satisfies", c.ver, c.constraint)
		if code != c.code || out != c.word+"\n" {
			t.Errorf("satisfies %s %q: out=%q code=%d, want %q/%d", c.ver, c.constraint, out, code, c.word, c.code)
		}
	}
}

func TestIntegrationSatisfiesMalformed(t *testing.T) {
	if _, _, code := run(t, "", "satisfies", "v1.0.0", ">=bogus"); code != 2 {
		t.Errorf("malformed constraint: code=%d, want 2", code)
	}
}

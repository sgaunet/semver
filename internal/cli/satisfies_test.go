package cli_test

import "testing"

func TestSatisfiesCommand(t *testing.T) {
	cases := []struct {
		ver, constraint, word string
		code                  int
	}{
		{"v1.5.0", "^1.2.0", "true", 0},
		{"v2.0.0", ">=1.2.0 <2.0.0", "false", 10},
		{"v1.5.0", ">=1.2.0 <2.0.0", "true", 0},
		{"v1.5.0-rc.1", ">=1.2.0 <2.0.0", "false", 10}, // pre-release exclusion
	}
	for _, c := range cases {
		out, errb, code := run("satisfies", c.ver, c.constraint)
		if code != c.code {
			t.Errorf("satisfies %s %q: exit = %d, want %d (stderr=%q)", c.ver, c.constraint, code, c.code, errb)
		}
		if out != c.word+"\n" {
			t.Errorf("satisfies %s %q: stdout = %q, want %q", c.ver, c.constraint, out, c.word)
		}
	}
}

func TestSatisfiesMalformedConstraint(t *testing.T) {
	if _, _, code := run("satisfies", "v1.0.0", ">=bogus"); code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}

func TestSatisfiesInvalidVersion(t *testing.T) {
	if _, _, code := run("satisfies", "bad", "^1.0.0"); code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}

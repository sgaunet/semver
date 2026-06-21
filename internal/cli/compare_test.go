package cli_test

import "testing"

func TestCompareCommand(t *testing.T) {
	cases := []struct {
		a, b string
		word string
		code int
	}{
		{"v1.0.0", "v1.2.0", "lower", 10},
		{"v2.0.0", "v1.9.9", "higher", 11},
		{"v1.0.0", "v1.0.0", "equal", 0},
		{"v1.0.0-rc.1", "v1.0.0", "lower", 10},
		{"v1.0.0+build.5", "v1.0.0+build.9", "equal", 0},
	}
	for _, c := range cases {
		out, errb, code := run("compare", c.a, c.b)
		if code != c.code {
			t.Errorf("compare %s %s: exit = %d, want %d", c.a, c.b, code, c.code)
		}
		if out != c.word+"\n" {
			t.Errorf("compare %s %s: stdout = %q, want %q", c.a, c.b, out, c.word)
		}
		if errb != "" {
			t.Errorf("compare %s %s: stderr = %q, want empty", c.a, c.b, errb)
		}
	}
}

func TestCompareInvalid(t *testing.T) {
	if _, _, code := run("compare", "v1.0.0", "bad"); code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
	if _, _, code := run("compare", "v1.0.0"); code != 2 {
		t.Errorf("one arg: exit = %d, want 2", code)
	}
}

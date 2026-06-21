package cli_test

import (
	"strings"
	"testing"
)

func TestBumpCommands(t *testing.T) {
	cases := []struct {
		args []string
		want string
	}{
		{[]string{"patch", "v1.0.0"}, "v1.0.1\n"},
		{[]string{"minor", "v1.2.3"}, "v1.3.0\n"},
		{[]string{"major", "v1.2.3"}, "v2.0.0\n"},
		{[]string{"patch", "1.0.0"}, "1.0.1\n"},
		{[]string{"patch", "v1.2.0-rc.1"}, "v1.2.0\n"},
	}
	for _, c := range cases {
		out, errb, code := run(c.args...)
		if code != 0 {
			t.Errorf("%v: exit = %d, want 0 (stderr=%q)", c.args, code, errb)
		}
		if out != c.want {
			t.Errorf("%v: stdout = %q, want %q", c.args, out, c.want)
		}
		if errb != "" {
			t.Errorf("%v: stderr = %q, want empty", c.args, errb)
		}
	}
}

func TestBumpInvalidVersion(t *testing.T) {
	out, errb, code := run("patch", "not-a-version")
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
	if out != "" {
		t.Errorf("stdout = %q, want empty on error", out)
	}
	if errb == "" {
		t.Error("want error message on stderr")
	}
}

func TestBumpMissingArg(t *testing.T) {
	if _, _, code := run("patch"); code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}

func TestBumpJSON(t *testing.T) {
	out, _, code := run("patch", "v1.0.0", "--output=json")
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	for _, want := range []string{`"input":"v1.0.0"`, `"operation":"patch"`, `"result":"v1.0.1"`} {
		if !strings.Contains(out, want) {
			t.Errorf("json %q missing %q", out, want)
		}
	}
}

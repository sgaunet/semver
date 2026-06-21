package cli_test

import "testing"

func TestPrereleaseStartAndIncrement(t *testing.T) {
	cases := []struct {
		args []string
		want string
	}{
		{[]string{"prerelease", "v1.0.0", "--pre", "rc", "--bump", "minor"}, "v1.1.0-rc.1\n"},
		{[]string{"prerelease", "v1.0.0", "--pre", "beta"}, "v1.0.1-beta.1\n"},
		{[]string{"prerelease", "v1.1.0-rc.1"}, "v1.1.0-rc.2\n"},
	}
	for _, c := range cases {
		out, errb, code := run(c.args...)
		if code != 0 {
			t.Errorf("%v: exit = %d (stderr=%q)", c.args, code, errb)
		}
		if out != c.want {
			t.Errorf("%v: stdout = %q, want %q", c.args, out, c.want)
		}
	}
}

func TestPrereleaseIncrementOnStableIsError(t *testing.T) {
	if _, _, code := run("prerelease", "v1.0.0"); code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}

func TestPrereleaseInvalidBump(t *testing.T) {
	if _, _, code := run("prerelease", "v1.0.0", "--pre", "rc", "--bump", "huge"); code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}

func TestReleaseCommand(t *testing.T) {
	out, _, code := run("release", "v1.1.0-rc.1")
	if code != 0 || out != "v1.1.0\n" {
		t.Errorf("release: out=%q code=%d", out, code)
	}
}

func TestReleaseOnStableIsError(t *testing.T) {
	if _, _, code := run("release", "v1.0.0"); code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}

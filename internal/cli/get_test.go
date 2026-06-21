package cli_test

import "testing"

func TestGetComponents(t *testing.T) {
	cases := []struct {
		comp, ver, want string
		code            int
	}{
		{"major", "v2.5.7-rc.3", "2\n", 0},
		{"minor", "v2.5.7-rc.3", "5\n", 0},
		{"patch", "v2.5.7-rc.3", "7\n", 0},
		{"prerelease", "v2.5.7-rc.3", "rc.3\n", 0},
		{"prerelease", "v2.5.7", "\n", 10}, // absent
		{"build", "v2.5.7+exp.1", "exp.1\n", 0},
		{"build", "v2.5.7", "\n", 10}, // absent
	}
	for _, c := range cases {
		out, _, code := run("get", c.comp, c.ver)
		if code != c.code {
			t.Errorf("get %s %s: exit = %d, want %d", c.comp, c.ver, code, c.code)
		}
		if out != c.want {
			t.Errorf("get %s %s: stdout = %q, want %q", c.comp, c.ver, out, c.want)
		}
	}
}

func TestGetUnknownComponent(t *testing.T) {
	if _, _, code := run("get", "flavor", "v1.0.0"); code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}

func TestGetInvalidVersion(t *testing.T) {
	if _, _, code := run("get", "major", "bad"); code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
}

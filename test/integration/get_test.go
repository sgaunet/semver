package integration

import "testing"

func TestIntegrationGet(t *testing.T) {
	cases := []struct {
		comp, ver, out string
		code           int
	}{
		{"major", "v2.5.7-rc.3", "2\n", 0},
		{"prerelease", "v2.5.7-rc.3", "rc.3\n", 0},
		{"prerelease", "v2.5.7", "\n", 10},
	}
	for _, c := range cases {
		out, _, code := run(t, "", "get", c.comp, c.ver)
		if code != c.code || out != c.out {
			t.Errorf("get %s %s: out=%q code=%d, want %q/%d", c.comp, c.ver, out, code, c.out, c.code)
		}
	}
}

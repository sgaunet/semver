package integration

import "testing"

func TestIntegrationBump(t *testing.T) {
	cases := []struct {
		args []string
		out  string
	}{
		{[]string{"patch", "v1.0.0"}, "v1.0.1\n"},
		{[]string{"minor", "v1.2.3"}, "v1.3.0\n"},
		{[]string{"major", "v1.2.3"}, "v2.0.0\n"},
		{[]string{"patch", "1.0.0"}, "1.0.1\n"},
	}
	for _, c := range cases {
		out, errb, code := run(t, "", c.args...)
		if code != 0 || out != c.out || errb != "" {
			t.Errorf("%v: out=%q err=%q code=%d, want out=%q code=0", c.args, out, errb, code, c.out)
		}
	}
}

func TestIntegrationBumpInvalid(t *testing.T) {
	out, errb, code := run(t, "", "patch", "1.0")
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
	if out != "" {
		t.Errorf("stdout = %q, want empty (data stream stays clean on error)", out)
	}
	if errb == "" {
		t.Error("want diagnostic on stderr")
	}
}

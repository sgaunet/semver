package integration

import "testing"

func TestIntegrationSortArgs(t *testing.T) {
	out, _, code := run(t, "", "sort", "v1.0.0", "v1.0.0-rc.1", "v2.0.0", "v1.2.0")
	want := "v1.0.0-rc.1\nv1.0.0\nv1.2.0\nv2.0.0\n"
	if code != 0 || out != want {
		t.Errorf("sort args: out=%q code=%d, want %q", out, code, want)
	}
}

func TestIntegrationSortStdin(t *testing.T) {
	out, _, code := run(t, "v2.0.0\nv1.0.0\nv1.2.0\n", "sort")
	want := "v1.0.0\nv1.2.0\nv2.0.0\n"
	if code != 0 || out != want {
		t.Errorf("sort stdin: out=%q code=%d, want %q", out, code, want)
	}
}

func TestIntegrationSortInvalidEntry(t *testing.T) {
	out, errb, code := run(t, "", "sort", "v1.0.0", "oops")
	if code != 2 || out != "" || errb == "" {
		t.Errorf("sort invalid: out=%q err=%q code=%d, want empty stdout + stderr + code 2", out, errb, code)
	}
}

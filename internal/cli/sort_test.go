package cli_test

import "testing"

func TestSortArgsAscending(t *testing.T) {
	out, _, code := run("sort", "v1.0.0", "v1.0.0-rc.1", "v2.0.0", "v1.2.0")
	want := "v1.0.0-rc.1\nv1.0.0\nv1.2.0\nv2.0.0\n"
	if code != 0 || out != want {
		t.Errorf("sort = %q (code %d), want %q", out, code, want)
	}
}

func TestSortDescending(t *testing.T) {
	out, _, code := run("sort", "--desc", "v1.0.0", "v2.0.0", "v1.5.0")
	want := "v2.0.0\nv1.5.0\nv1.0.0\n"
	if code != 0 || out != want {
		t.Errorf("sort --desc = %q (code %d), want %q", out, code, want)
	}
}

func TestSortStdin(t *testing.T) {
	out, _, code := runStdin("v2.0.0\nv1.0.0\nv1.2.0\n", "sort")
	want := "v1.0.0\nv1.2.0\nv2.0.0\n"
	if code != 0 || out != want {
		t.Errorf("sort stdin = %q (code %d), want %q", out, code, want)
	}
}

func TestSortStableForEqualPrecedence(t *testing.T) {
	// Build metadata does not affect precedence, so input order is preserved.
	out, _, code := run("sort", "v1.0.0+b", "v1.0.0+a")
	want := "v1.0.0+b\nv1.0.0+a\n"
	if code != 0 || out != want {
		t.Errorf("stable sort = %q (code %d), want %q", out, code, want)
	}
}

func TestSortInvalidEntryFails(t *testing.T) {
	out, errb, code := run("sort", "v1.0.0", "nope", "v2.0.0")
	if code != 2 {
		t.Errorf("exit = %d, want 2", code)
	}
	if out != "" {
		t.Errorf("stdout = %q, want empty on error", out)
	}
	if errb == "" {
		t.Error("want error naming the offending entry on stderr")
	}
}

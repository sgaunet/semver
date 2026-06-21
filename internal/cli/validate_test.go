package cli_test

import (
	"strings"
	"testing"
)

func TestValidateValid(t *testing.T) {
	out, errb, code := run("validate", "v1.2.3-rc.1+build.7")
	if code != 0 {
		t.Errorf("exit = %d, want 0", code)
	}
	if out != "valid\n" {
		t.Errorf("stdout = %q, want valid", out)
	}
	if errb != "" {
		t.Errorf("stderr = %q, want empty", errb)
	}
}

func TestValidateInvalid(t *testing.T) {
	for _, in := range []string{"1.2", "01.2.3", ""} {
		out, _, code := run("validate", in)
		if code != 10 {
			t.Errorf("validate %q: exit = %d, want 10", in, code)
		}
		if out != "" {
			t.Errorf("validate %q: stdout = %q, want empty", in, out)
		}
	}
}

func TestValidateJSONComponents(t *testing.T) {
	out, _, code := run("validate", "v1.2.3-rc.1+build.7", "--output=json")
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	for _, want := range []string{`"valid":true`, `"major":1`, `"minor":2`, `"patch":3`, `"prerelease":"rc.1"`, `"build":"build.7"`} {
		if !strings.Contains(out, want) {
			t.Errorf("json %q missing %q", out, want)
		}
	}
}

func TestValidateJSONInvalid(t *testing.T) {
	out, _, code := run("validate", "1.2", "--output=json")
	if code != 10 {
		t.Errorf("exit = %d, want 10", code)
	}
	if !strings.Contains(out, `"valid":false`) {
		t.Errorf("json %q missing valid:false", out)
	}
}

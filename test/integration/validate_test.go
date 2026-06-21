package integration

import (
	"strings"
	"testing"
)

func TestIntegrationValidate(t *testing.T) {
	out, _, code := run(t, "", "validate", "v1.2.3-rc.1+build.7")
	if code != 0 || out != "valid\n" {
		t.Errorf("valid: out=%q code=%d", out, code)
	}

	out, errb, code := run(t, "", "validate", "1.2")
	if code != 10 {
		t.Errorf("invalid: code=%d, want 10", code)
	}
	if out != "" || errb == "" {
		t.Errorf("invalid: out=%q err=%q, want empty stdout + stderr message", out, errb)
	}
}

func TestIntegrationValidateJSON(t *testing.T) {
	out, _, code := run(t, "", "validate", "v1.2.3-rc.1+build.7", "--output=json")
	if code != 0 {
		t.Fatalf("code=%d", code)
	}
	if !strings.Contains(out, `"valid":true`) || !strings.Contains(out, `"prerelease":"rc.1"`) {
		t.Errorf("json = %q", out)
	}
}

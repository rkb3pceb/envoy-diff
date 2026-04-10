package lint

import (
	"strings"
	"testing"
)

func TestLint_NoFindings_CleanEnv(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost:5432/mydb",
		"APP_PORT":     "8080",
	}
	findings := Lint(env)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d: %v", len(findings), findings)
	}
}

func TestLint_LowercaseKey_ReturnsWarn(t *testing.T) {
	env := map[string]string{"app_port": "8080"}
	findings := Lint(env)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Severity != SeverityWarn {
		t.Errorf("expected warn, got %s", findings[0].Severity)
	}
	if !strings.Contains(findings[0].Message, "uppercase") {
		t.Errorf("unexpected message: %s", findings[0].Message)
	}
}

func TestLint_KeyWithSpace_ReturnsError(t *testing.T) {
	env := map[string]string{"APP PORT": "8080"}
	findings := Lint(env)
	var errs []Finding
	for _, f := range findings {
		if f.Severity == SeverityError {
			errs = append(errs, f)
		}
	}
	if len(errs) == 0 {
		t.Fatal("expected at least one error finding for key with space")
	}
}

func TestLint_EmptyValue_ReturnsWarn(t *testing.T) {
	env := map[string]string{"SECRET_KEY": ""}
	findings := Lint(env)
	if len(findings) == 0 {
		t.Fatal("expected finding for empty value")
	}
	for _, f := range findings {
		if f.Key == "SECRET_KEY" && strings.Contains(f.Message, "empty") {
			return
		}
	}
	t.Error("did not find expected empty-value warning")
}

func TestLint_PlaceholderValue_ReturnsWarn(t *testing.T) {
	cases := []string{"CHANGEME", "your_secret", "<TOKEN>", "TODO"}
	for _, v := range cases {
		env := map[string]string{"API_KEY": v}
		findings := Lint(env)
		found := false
		for _, f := range findings {
			if strings.Contains(f.Message, "placeholder") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected placeholder warning for value %q", v)
		}
	}
}

func TestFinding_String_Format(t *testing.T) {
	f := Finding{Key: "FOO", Message: "some issue", Severity: SeverityError}
	s := f.String()
	if !strings.Contains(s, "error") || !strings.Contains(s, "FOO") || !strings.Contains(s, "some issue") {
		t.Errorf("unexpected String() output: %s", s)
	}
}

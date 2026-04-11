package template_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envoy-diff/internal/template"
)

func TestRender_PlainText_Unchanged(t *testing.T) {
	src := []byte("KEY=value\nOTHER=123\n")
	out, err := template.Render(src, template.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(src) {
		t.Errorf("expected %q, got %q", src, out)
	}
}

func TestRender_SubstitutesVar(t *testing.T) {
	src := []byte("APP_ENV={{ .DEPLOY_ENV }}\n")
	opts := template.DefaultOptions()
	opts.Vars["DEPLOY_ENV"] = "production"

	out, err := template.Render(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "APP_ENV=production\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRender_DefaultFunc(t *testing.T) {
	src := []byte(`DB_HOST={{ default "localhost" .DB_HOST }}`)
	out, err := template.Render(src, template.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "DB_HOST=localhost" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRender_MissingVar_ReturnsError(t *testing.T) {
	src := []byte("X={{ .UNDEFINED_XYZ_VAR_12345 }}\n")
	opts := template.DefaultOptions()
	opts.AllowMissing = false

	_, err := template.Render(src, opts)
	if err == nil {
		t.Fatal("expected error for missing variable")
	}
}

func TestRender_AllowMissing_EmptyString(t *testing.T) {
	src := []byte("X={{ .UNDEFINED_XYZ_VAR_12345 }}\n")
	opts := template.DefaultOptions()
	opts.AllowMissing = true

	out, err := template.Render(src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "X=\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRenderFile_ReadsAndRenders(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "env.tmpl")
	_ = os.WriteFile(path, []byte("REGION={{ .REGION }}\n"), 0o644)

	opts := template.DefaultOptions()
	opts.Vars["REGION"] = "eu-west-1"

	out, err := template.RenderFile(path, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "REGION=eu-west-1\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRenderFile_MissingFile_ReturnsError(t *testing.T) {
	_, err := template.RenderFile("/nonexistent/path/env.tmpl", template.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

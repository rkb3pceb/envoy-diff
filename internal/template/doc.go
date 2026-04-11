// Package template renders .env files that contain Go template directives.
//
// It is useful when a single template file needs to produce environment
// configurations for multiple deployment targets (e.g. staging vs production)
// without duplicating the bulk of the file.
//
// # Basic usage
//
//	// Render an in-memory template with extra variables.
//	out, err := template.Render(src, template.Options{
//		Vars: map[string]string{"ENV": "staging"},
//	})
//
//	// Render a template file on disk.
//	out, err := template.RenderFile("deploy/staging.env.tmpl", opts)
//
// # Template syntax
//
// Templates follow standard Go text/template syntax. The context object is a
// map[string]string populated from the OS environment merged with any Vars
// supplied via Options (Vars take precedence).
//
// Built-in helper functions:
//
//	{{ default "fallback" .MY_VAR }}  — return fallback when MY_VAR is empty
//	{{ upper .MY_VAR }}               — convert to upper-case
//	{{ lower .MY_VAR }}               — convert to lower-case
package template

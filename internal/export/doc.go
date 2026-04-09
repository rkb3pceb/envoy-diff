// Package export provides exporters that serialize diff results
// into external file formats for reporting and archival purposes.
//
// Supported formats:
//
//	"csv"      — comma-separated values, suitable for spreadsheet tools
//	"markdown" — GitHub-flavored Markdown table, suitable for PR comments
//
// Usage:
//
//	e, err := export.New(export.FormatMarkdown)
//	if err != nil {
//		log.Fatal(err)
//	}
//	opts := export.Options{
//		Format:    export.FormatMarkdown,
//		OldFile:   "staging.env",
//		NewFile:   "production.env",
//		Timestamp: time.Now(),
//	}
//	if err := e.Write(os.Stdout, changes, opts); err != nil {
//		log.Fatal(err)
//	}
package export

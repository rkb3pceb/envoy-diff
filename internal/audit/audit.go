// Package audit provides security and compliance auditing for environment variable diffs.
package audit

import (
	"strings"

	"github.com/user/envoy-diff/internal/diff"
)

// Severity represents the risk level of an audit finding.
type Severity string

const (
	SeverityHigh   Severity = "HIGH"
	SeverityMedium Severity = "MEDIUM"
	SeverityLow    Severity = "LOW"
)

// Finding represents a single audit finding for a changed environment variable.
type Finding struct {
	Key      string
	Severity Severity
	Reason   string
	Change   diff.Change
}

// sensitivePatterns are substrings that indicate a sensitive variable.
var sensitivePatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
	"DSN",
	"DATABASE_URL",
}

// Audit inspects a list of diff changes and returns findings for sensitive
// or potentially dangerous variable modifications.
func Audit(changes []diff.Change) []Finding {
	var findings []Finding

	for _, c := range changes {
		upper := strings.ToUpper(c.Key)

		if c.Type == diff.Removed {
			findings = append(findings, Finding{
				Key:      c.Key,
				Severity: SeverityMedium,
				Reason:   "variable removed from deployment config",
				Change:   c,
			})
			continue
		}

		for _, pattern := range sensitivePatterns {
			if strings.Contains(upper, pattern) {
				severity := SeverityHigh
				if c.Type == diff.Added {
					severity = SeverityMedium
				}
				findings = append(findings, Finding{
					Key:      c.Key,
					Severity: severity,
					Reason:   "sensitive variable " + strings.ToLower(string(c.Type)),
					Change:   c,
				})
				break
			}
		}
	}

	return findings
}

// Package notify implements pluggable notification backends for envoy-diff.
//
// When a diff run produces audit findings that meet or exceed a configured
// severity threshold, the Dispatcher fans out the event to every registered
// Notifier.
//
// Built-in notifiers:
//
//	StdoutNotifier  — writes a plain-text summary to any io.Writer.
//	WebhookNotifier — POSTs a JSON payload to an HTTP endpoint.
//
// Custom notifiers can be added by implementing the Notifier interface:
//
//	type Notifier interface {
//	    Send(e Event) error
//	}
//
// Example:
//
//	d := notify.New(notify.LevelHigh,
//	    notify.NewStdoutNotifier(os.Stderr),
//	    notify.NewWebhookNotifier("https://hooks.example.com/envoy"),
//	)
//	_ = d.Dispatch(event)
package notify

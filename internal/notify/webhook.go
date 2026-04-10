package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookNotifier sends a JSON payload to an HTTP endpoint.
type WebhookNotifier struct {
	URL    string
	client *http.Client
}

// NewWebhookNotifier returns a WebhookNotifier targeting url.
func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{
		URL: url,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type webhookPayload struct {
	Summary  string `json:"summary"`
	Changes  int    `json:"changes"`
	Findings int    `json:"findings"`
}

// Send marshals the event and POSTs it to the configured URL.
func (w *WebhookNotifier) Send(e Event) error {
	payload := webhookPayload{
		Summary:  e.Summary,
		Changes:  len(e.Changes),
		Findings: len(e.Findings),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal: %w", err)
	}
	resp, err := w.client.Post(w.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}

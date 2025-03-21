package alert

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type webhookSender struct {
	httpClient *http.Client
	url        string
	verbose    bool
	failedOnly bool
}

func NewWebhookSender(httpClient *http.Client, url string, verbose bool, failedOnly bool) *webhookSender {
	return &webhookSender{
		httpClient,
		url,
		verbose,
		failedOnly,
	}
}

func (s *webhookSender) send(m Results, failedCount int) error {
	// If there are no failed sockets and we only wish to be notified when they fail, there is nothing to do
	if failedCount == 0 && s.failedOnly {
		log.Printf("%T: no failed sockets and failedOnly == true, nothing will be sent", s)
		return nil
	}

	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	if s.verbose {
		log.Printf("Prepared webhook data: %v", string(jsonData))
	}

	resp, err := s.httpClient.Post(s.url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if s.verbose {
		log.Printf("Webhook notification sent. Webhook URL: %s", s.url)
		log.Printf("Received response from webhook URL. Status: %s. Body: %s", resp.Status, string(body))
	}

	return nil
}

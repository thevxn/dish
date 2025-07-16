package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/logger"
)

type webhookSender struct {
	httpClient    HTTPClient
	url           string
	verbose       bool
	notifySuccess bool
	logger        logger.Logger
}

func NewWebhookSender(httpClient HTTPClient, config *config.Config, logger logger.Logger) (*webhookSender, error) {
	parsedURL, err := parseAndValidateURL(config.WebhookURL, nil)
	if err != nil {
		return nil, err
	}

	return &webhookSender{
		httpClient:    httpClient,
		url:           parsedURL.String(),
		verbose:       config.Verbose,
		notifySuccess: config.MachineNotifySuccess,
		logger:        logger,
	}, nil
}

func (s *webhookSender) send(m *Results, failedCount int) error {
	// If no checks failed and success should not be notified, there is nothing to send
	if failedCount == 0 && !s.notifySuccess {
		s.logger.Debug("no sockets failed, nothing will be sent to webhook")

		return nil
	}

	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(jsonData)

	s.logger.Debugf("prepared webhook data: %s", string(jsonData))

	res, err := handleSubmit(s.httpClient, http.MethodPost, s.url, bodyReader)
	if err != nil {
		return fmt.Errorf("error pushing results to webhook: %w", err)
	}

	err = handleRead(res, s.logger)
	if err != nil {
		return err
	}

	s.logger.Info("results pushed to webhook")

	return nil
}

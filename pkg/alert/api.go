package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/logger"
)

type apiSender struct {
	httpClient    HTTPClient
	url           string
	headerName    string
	headerValue   string
	verbose       bool
	notifySuccess bool
	logger        logger.Logger
}

func NewAPISender(httpClient HTTPClient, config *config.Config, logger logger.Logger) (*apiSender, error) {
	parsedURL, err := parseAndValidateURL(config.ApiURL, nil)
	if err != nil {
		return nil, err
	}

	return &apiSender{
		httpClient:    httpClient,
		url:           parsedURL.String(),
		headerName:    config.ApiHeaderName,
		headerValue:   config.ApiHeaderValue,
		verbose:       config.Verbose,
		notifySuccess: config.MachineNotifySuccess,
		logger:        logger,
	}, nil
}

func (s *apiSender) send(m *Results, failedCount int) error {
	// If no checks failed and success should not be notified, there is nothing to send
	if failedCount == 0 && !s.notifySuccess {
		s.logger.Debug("no sockets failed, nothing will be sent to remote API")

		return nil
	}

	jsonData, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	bodyReader := bytes.NewReader(jsonData)

	s.logger.Debugf("prepared remote API data: %s", string(jsonData))

	// If custom header & value is provided (mostly used for auth purposes), include it in the request
	opts := []func(*submitOptions){}
	if s.headerName != "" && s.headerValue != "" {
		opts = append(opts, withHeader(s.headerName, s.headerValue))
	}

	res, err := handleSubmit(s.httpClient, http.MethodPost, s.url, bodyReader, opts...)
	if err != nil {
		return fmt.Errorf("error pushing results to remote API: %w", err)
	}

	err = handleRead(res, s.logger)
	if err != nil {
		return err
	}

	s.logger.Info("results pushed to remote API")

	return nil
}

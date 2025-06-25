package alert

import (
	"reflect"
	"testing"

	"go.vxn.dev/dish/pkg/config"
	testhelpers "go.vxn.dev/dish/pkg/testdata"
)

func TestNewWebhookSender(t *testing.T) {
	mockHTTPClient := &testhelpers.SuccessStatusHTTPClient{}
	mockLogger := &testhelpers.MockLogger{}

	url := "https://abc123.xyz.com"
	notifySuccess := false
	verbose := false

	expected := &webhookSender{
		httpClient:    mockHTTPClient,
		url:           url,
		notifySuccess: notifySuccess,
		verbose:       verbose,
		logger:        mockLogger,
	}

	cfg := &config.Config{
		WebhookURL:           url,
		Verbose:              verbose,
		MachineNotifySuccess: notifySuccess,
	}
	actual, _ := NewWebhookSender(mockHTTPClient, cfg, mockLogger)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestSend_Webhook(t *testing.T) {
	url := "https://abc123.xyz.com"

	successResults := Results{
		Map: map[string]bool{
			"test": true,
		},
	}
	failedResults := Results{
		Map: map[string]bool{
			"test": false,
		},
	}
	mixedResults := Results{
		Map: map[string]bool{
			"test1": true,
			"test2": false,
		},
	}

	newConfig := func(url string, notifySuccess, verbose bool) *config.Config {
		return &config.Config{
			WebhookURL:           url,
			Verbose:              verbose,
			MachineNotifySuccess: notifySuccess,
		}
	}

	tests := []struct {
		name          string
		client        HTTPClient
		results       Results
		failedCount   int
		notifySuccess bool
		verbose       bool
		wantErr       bool
	}{
		{
			name:          "Failed Sockets",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			results:       failedResults,
			failedCount:   1,
			notifySuccess: false,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "Failed Sockets - Verbose",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			results:       failedResults,
			failedCount:   1,
			notifySuccess: false,
			verbose:       true,
			wantErr:       false,
		},
		{
			name:          "No Failed Sockets With notifySuccess",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			results:       successResults,
			failedCount:   0,
			notifySuccess: true,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "No Failed Sockets Without notifySuccess",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			results:       successResults,
			failedCount:   0,
			notifySuccess: false,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "No Failed Sockets Without notifySuccess - Verbose",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			results:       successResults,
			failedCount:   0,
			notifySuccess: false,
			verbose:       true,
			wantErr:       false,
		},
		{
			name:          "Mixed Results With notifySuccess",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			results:       mixedResults,
			failedCount:   1,
			notifySuccess: true,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "Mixed Results Without notifySuccess",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			results:       mixedResults,
			failedCount:   1,
			notifySuccess: false,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "Network Error When Pushing to Webhook",
			client:        &testhelpers.FailureHTTPClient{},
			results:       failedResults,
			failedCount:   1,
			notifySuccess: false,
			verbose:       false,
			wantErr:       true,
		},
		{
			name:          "Unexpected Response Code From Webhook",
			client:        &testhelpers.ErrorStatusHTTPClient{},
			results:       failedResults,
			failedCount:   1,
			notifySuccess: false,
			verbose:       false,
			wantErr:       true,
		},
		{
			name:          "Error Reading Response Body From Webhook",
			client:        &testhelpers.InvalidResponseBodyHTTPClient{},
			results:       failedResults,
			failedCount:   1,
			notifySuccess: false,
			verbose:       true,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newConfig(url, tt.verbose, tt.notifySuccess)
			sender, err := NewWebhookSender(tt.client, cfg, &testhelpers.MockLogger{})
			if err != nil {
				t.Fatalf("failed to create Webhook sender instance: %v", err)
			}

			err = sender.send(&tt.results, tt.failedCount)
			if tt.wantErr != (err != nil) {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

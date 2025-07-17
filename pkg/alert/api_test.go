package alert

import (
	"reflect"
	"testing"

	"go.vxn.dev/dish/pkg/config"
)

func TestNewAPISender(t *testing.T) {
	mockHTTPClient := &SuccessStatusHTTPClient{}

	headerName := "X-Api-Key"
	headerValue := "abc123"
	notifySuccess := false
	verbose := false

	expected := &apiSender{
		httpClient:    mockHTTPClient,
		url:           pushgatewayURL,
		headerName:    headerName,
		headerValue:   headerValue,
		notifySuccess: notifySuccess,
		verbose:       verbose,
		logger:        &MockLogger{},
	}

	cfg := &config.Config{
		ApiURL:               pushgatewayURL,
		ApiHeaderName:        headerName,
		ApiHeaderValue:       headerValue,
		MachineNotifySuccess: notifySuccess,
		Verbose:              verbose,
	}

	actual, _ := NewAPISender(mockHTTPClient, cfg, &MockLogger{})

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestSend_API(t *testing.T) {
	headerName := "X-Api-Key"
	headerValue := "abc123"

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

	newConfig := func(headerName, headerValue string, notifySuccess, verbose bool) *config.Config {
		return &config.Config{
			ApiURL:               pushgatewayURL,
			MachineNotifySuccess: notifySuccess,
			Verbose:              verbose,
			ApiHeaderName:        headerName,
			ApiHeaderValue:       headerValue,
		}
	}

	tests := []struct {
		name          string
		client        HTTPClient
		results       Results
		failedCount   int
		notifySuccess bool
		headerName    string
		headerValue   string
		verbose       bool
		wantErr       bool
	}{
		{
			name:          "Failed Sockets",
			client:        &SuccessStatusHTTPClient{},
			results:       failedResults,
			failedCount:   1,
			notifySuccess: false,
			headerName:    headerName,
			headerValue:   headerValue,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "Failed Sockets - Verbose",
			client:        &SuccessStatusHTTPClient{},
			results:       failedResults,
			failedCount:   1,
			notifySuccess: false,
			headerName:    headerName,
			headerValue:   headerValue,
			verbose:       true,
			wantErr:       false,
		},
		{
			name:          "No Failed Sockets With notifySuccess",
			client:        &SuccessStatusHTTPClient{},
			results:       successResults,
			failedCount:   0,
			notifySuccess: true,
			headerName:    headerName,
			headerValue:   headerValue,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "No Failed Sockets Without notifySuccess",
			client:        &SuccessStatusHTTPClient{},
			results:       successResults,
			failedCount:   0,
			notifySuccess: false,
			headerName:    headerName,
			headerValue:   headerValue,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "No Failed Sockets Without notifySuccess - Verbose",
			client:        &SuccessStatusHTTPClient{},
			results:       successResults,
			failedCount:   0,
			notifySuccess: false,
			headerName:    headerName,
			headerValue:   headerValue,
			verbose:       true,
			wantErr:       false,
		},
		{
			name:          "Mixed Results With notifySuccess",
			client:        &SuccessStatusHTTPClient{},
			results:       mixedResults,
			failedCount:   1,
			notifySuccess: true,
			headerName:    "",
			headerValue:   "",
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "Mixed Results Without notifySuccess",
			client:        &SuccessStatusHTTPClient{},
			results:       mixedResults,
			failedCount:   1,
			notifySuccess: false,
			headerName:    "",
			headerValue:   "",
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "No Custom Header",
			client:        &SuccessStatusHTTPClient{},
			results:       failedResults,
			failedCount:   1,
			notifySuccess: false,
			headerName:    "",
			headerValue:   "",
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "Network Error When Pushing to Remote API",
			client:        &FailureHTTPClient{},
			results:       failedResults,
			failedCount:   1,
			notifySuccess: false,
			headerName:    headerName,
			headerValue:   headerValue,
			verbose:       false,
			wantErr:       true,
		},
		{
			name:          "Unexpected Response Code From Remote API",
			client:        &ErrorStatusHTTPClient{},
			results:       failedResults,
			failedCount:   1,
			notifySuccess: false,
			headerName:    headerName,
			headerValue:   headerValue,
			verbose:       false,
			wantErr:       true,
		},
		{
			name:          "Error Reading Response Body From Remote API",
			client:        &InvalidResponseBodyHTTPClient{},
			results:       failedResults,
			failedCount:   1,
			notifySuccess: false,
			headerName:    headerName,
			headerValue:   headerValue,
			verbose:       true,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newConfig(tt.headerName, tt.headerValue, tt.notifySuccess, tt.verbose)
			sender, err := NewAPISender(tt.client, cfg, &MockLogger{})
			if err != nil {
				t.Fatalf("failed to create API sender instance: %v", err)
			}

			err = sender.send(&tt.results, tt.failedCount)
			if tt.wantErr != (err != nil) {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

package alert

import (
	"fmt"
	"testing"

	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/testhelpers"
)

func TestNewPushgatewaySender(t *testing.T) {
	mockHTTPClient := &testhelpers.SuccessStatusHTTPClient{}
	mockLogger := &testhelpers.MockLogger{}

	url := "https://abc123.xyz.com"
	instanceName := "test-instance"
	verbose := false
	notifySuccess := false

	expected := &pushgatewaySender{
		httpClient:    mockHTTPClient,
		url:           url,
		instanceName:  "test-instance",
		notifySuccess: notifySuccess,
		verbose:       verbose,
		logger:        mockLogger,
		// template will be compared based on its output, no need for it here
	}

	cfg := &config.Config{
		PushgatewayURL:       url,
		InstanceName:         instanceName,
		Verbose:              verbose,
		MachineNotifySuccess: notifySuccess,
	}

	actual, err := NewPushgatewaySender(mockHTTPClient, cfg, mockLogger)
	if err != nil {
		t.Fatalf("error creating a new Pushgateway sender instance: %v", err)
	}

	// Compare fields individually due to complex structs
	if expected.url != actual.url {
		t.Errorf("expected url: %s, got: %s", expected.url, actual.url)
	}
	if expected.instanceName != actual.instanceName {
		t.Errorf("expected instanceName: %s, got: %s", expected.instanceName, actual.instanceName)
	}
	if expected.verbose != actual.verbose {
		t.Errorf("expected verbose: %v, got: %v", expected.verbose, actual.verbose)
	}
	if expected.notifySuccess != actual.notifySuccess {
		t.Errorf("expected notifySuccess: %v, got: %v", expected.notifySuccess, actual.notifySuccess)
	}
	if fmt.Sprintf("%T", expected.httpClient) != fmt.Sprintf("%T", actual.httpClient) {
		t.Errorf("expected httpClient type: %T, got: %T", expected.httpClient, actual.httpClient)
	}
	if fmt.Sprintf("%T", expected.logger) != fmt.Sprintf("%T", actual.logger) {
		t.Errorf("expected logger type: %T, got: %T", expected.logger, actual.logger)
	}
}

func TestSend_Pushgateway(t *testing.T) {
	url := "https://abc123.xyz.com"
	instanceName := "test-instance"

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

	newConfig := func(url, instanceName string, notifySuccess, verbose bool) *config.Config {
		return &config.Config{
			PushgatewayURL:       url,
			InstanceName:         instanceName,
			Verbose:              verbose,
			MachineNotifySuccess: notifySuccess,
		}
	}

	tests := []struct {
		name          string
		client        HTTPClient
		results       Results
		failedCount   int
		instanceName  string
		notifySuccess bool
		verbose       bool
		wantErr       bool
	}{
		{
			name:          "Failed Sockets",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			results:       failedResults,
			failedCount:   1,
			instanceName:  instanceName,
			notifySuccess: false,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "Failed Sockets - Verbose",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			results:       failedResults,
			failedCount:   1,
			instanceName:  instanceName,
			notifySuccess: false,
			verbose:       true,
			wantErr:       false,
		},
		{
			name:          "No Failed Sockets With notifySuccess",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			results:       successResults,
			failedCount:   0,
			instanceName:  instanceName,
			notifySuccess: true,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "No Failed Sockets Without notifySuccess",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			results:       successResults,
			failedCount:   0,
			instanceName:  instanceName,
			notifySuccess: false,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "No Failed Sockets Without notifySuccess - Verbose",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			results:       successResults,
			failedCount:   0,
			instanceName:  instanceName,
			notifySuccess: false,
			verbose:       true,
			wantErr:       false,
		},
		{
			name:          "Mixed Results With notifySuccess",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			results:       mixedResults,
			failedCount:   1,
			instanceName:  instanceName,
			notifySuccess: true,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "Mixed Results Without notifySuccess",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			results:       mixedResults,
			failedCount:   1,
			instanceName:  instanceName,
			notifySuccess: false,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "Empty Instance Name",
			client:        &testhelpers.SuccessStatusHTTPClient{},
			results:       failedResults,
			failedCount:   1,
			instanceName:  "",
			notifySuccess: false,
			verbose:       false,
			wantErr:       false,
		},
		{
			name:          "Network Error When Pushing to Pushgateway",
			client:        &testhelpers.FailureHTTPClient{},
			results:       failedResults,
			failedCount:   1,
			instanceName:  instanceName,
			notifySuccess: false,
			verbose:       false,
			wantErr:       true,
		},
		{
			name:          "Unexpected Response Code From Pushgateway",
			client:        &testhelpers.ErrorStatusHTTPClient{},
			results:       failedResults,
			failedCount:   1,
			instanceName:  instanceName,
			notifySuccess: false,
			verbose:       false,
			wantErr:       true,
		},
		{
			name:          "Error Reading Response Body From Pushgateway",
			client:        &testhelpers.InvalidResponseBodyHTTPClient{},
			results:       failedResults,
			failedCount:   1,
			instanceName:  instanceName,
			notifySuccess: false,
			verbose:       true,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newConfig(url, tt.instanceName, tt.notifySuccess, tt.verbose)
			sender, err := NewPushgatewaySender(tt.client, cfg, &testhelpers.MockLogger{})
			if err != nil {
				t.Fatalf("failed to create Pushgateway sender instance: %v", err)
			}

			err = sender.send(&tt.results, tt.failedCount)
			if tt.wantErr != (err != nil) {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestCreateMessage(t *testing.T) {
	cfg := &config.Config{
		PushgatewayURL:       "https://abc123.xyz.com",
		InstanceName:         "test-instance",
		MachineNotifySuccess: false,
		Verbose:              false,
	}

	sender, err := NewPushgatewaySender(&testhelpers.SuccessStatusHTTPClient{}, cfg, &testhelpers.MockLogger{})
	if err != nil {
		t.Fatalf("failed to create Pushgateway sender instance: %v", err)
	}

	failedCount := 1

	expected := `
#HELP failed sockets registered by dish
#TYPE dish_failed_count counter
dish_failed_count 1

`

	actual, err := sender.createMessage(failedCount)
	if err != nil {
		t.Errorf("error creating Pushgateway message: %v", err)
	}

	if expected != actual {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

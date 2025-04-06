package alert

import (
	"reflect"
	"testing"

	"go.vxn.dev/dish/pkg/testhelpers"
)

func TestNewAPISender(t *testing.T) {
	mockHTTPClient := &testhelpers.SuccessStatusHTTPClient{}

	url := "https://abc123.xyz.com"
	verbose := false
	notifySuccess := false
	headerName := "X-Api-Key"
	headerValue := "abc123"

	expected := &apiSender{
		httpClient:    mockHTTPClient,
		verbose:       verbose,
		notifySuccess: notifySuccess,
		headerName:    headerName,
		headerValue:   headerValue,
		url:           url,
	}

	actual, _ := NewAPISender(mockHTTPClient, url, headerName, headerValue, verbose, notifySuccess)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestSend_API(t *testing.T) {
	url := "https://abc123.xyz.com"
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
			client:        &testhelpers.SuccessStatusHTTPClient{},
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
			client:        &testhelpers.SuccessStatusHTTPClient{},
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
			client:        &testhelpers.SuccessStatusHTTPClient{},
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
			client:        &testhelpers.SuccessStatusHTTPClient{},
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
			client:        &testhelpers.SuccessStatusHTTPClient{},
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
			client:        &testhelpers.SuccessStatusHTTPClient{},
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
			client:        &testhelpers.SuccessStatusHTTPClient{},
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
			client:        &testhelpers.SuccessStatusHTTPClient{},
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
			client:        &testhelpers.FailureHTTPClient{},
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
			client:        &testhelpers.ErrorStatusHTTPClient{},
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
			client:        &testhelpers.InvalidResponseBodyHTTPClient{},
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
			sender, err := NewAPISender(tt.client, url, tt.headerName, tt.headerValue, tt.verbose, tt.notifySuccess)
			if err != nil {
				t.Fatalf("failed to create API sender instance: %v", err)
			}

			err = sender.send(tt.results, tt.failedCount)
			if tt.wantErr != (err != nil) {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

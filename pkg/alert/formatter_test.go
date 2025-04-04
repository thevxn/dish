package alert

import (
	"errors"
	"fmt"
	"testing"

	"go.vxn.dev/dish/pkg/socket"
)

func TestFormatMessengerText(t *testing.T) {
	tests := []struct {
		name         string
		result       socket.Result
		expectedText string
	}{
		{
			name: "Passed TCP Check",
			result: socket.Result{
				Socket: socket.Socket{
					ID:   "test_socket",
					Name: "test socket",
					Host: "192.168.0.1",
					Port: 123,
				},
				Passed: true,
				Error:  nil,
			},
			expectedText: "• 192.168.0.1:123 -- success ✅\n",
		},
		{
			name: "Passed HTTP Check",
			result: socket.Result{
				Socket: socket.Socket{
					ID:                "test_socket",
					Name:              "test socket",
					Host:              "https://test.testdomain.xyz",
					Port:              80,
					ExpectedHTTPCodes: []int{200},
					PathHTTP:          "/",
				},
				Passed: true,
				Error:  nil,
			},
			expectedText: "• https://test.testdomain.xyz:80/ -- success ✅\n",
		},
		{
			name: "Failed TCP Check",
			result: socket.Result{
				Socket: socket.Socket{
					ID:   "test_socket",
					Name: "test socket",
					Host: "192.168.0.1",
					Port: 123,
				},
				Passed: false,
				Error:  errors.New("error message"),
			},
			expectedText: "• 192.168.0.1:123 -- failed ❌ -- error message\n",
		},
		{
			name: "Failed HTTP Check with Error",
			result: socket.Result{
				Socket: socket.Socket{
					ID:                "test_socket",
					Name:              "test socket",
					Host:              "https://test.testdomain.xyz",
					Port:              80,
					ExpectedHTTPCodes: []int{200},
					PathHTTP:          "/",
				},
				Passed: false,
				Error:  errors.New("error message"),
			},
			expectedText: "• https://test.testdomain.xyz:80/ -- failed ❌ -- error message\n",
		},
		{
			name: "Failed HTTP Check with Unexpected Response Code",
			result: socket.Result{
				Socket: socket.Socket{
					ID:                "test_socket",
					Name:              "test socket",
					Host:              "https://test.testdomain.xyz",
					Port:              80,
					ExpectedHTTPCodes: []int{200},
					PathHTTP:          "/",
				},
				ResponseCode: 500,
				Passed:       false,
				Error:        fmt.Errorf("expected codes: %v, got %d", []int{200}, 500),
			},
			expectedText: "• https://test.testdomain.xyz:80/ -- failed ❌ -- expected codes: [200], got 500\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualText := FormatMessengerText(tt.result)

			if actualText != tt.expectedText {
				t.Errorf("expected %s, got %s", tt.expectedText, actualText)
			}
		})
	}
}

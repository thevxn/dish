package netrunner

import (
	"context"
	"flag"
	"net/http"
	"reflect"
	"runtime"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/logger"
	"go.vxn.dev/dish/pkg/socket"
)

// TestChecksum tests the checksum calculation function
func TestChecksum(t *testing.T) {
	if runtime.GOOS == "windows" {
		expected := 0
		actual := checksum([]byte{})

		if expected != int(actual) {
			t.Errorf("unexpected windows checksum. expected: %d, got: %d", expected, actual)
		}
		return
	}

	tests := []struct {
		name     string
		input    []byte
		expected uint16
	}{
		{
			name:     "empty slice",
			input:    []byte{},
			expected: 0xFFFF,
		},
		{
			name:     "single byte",
			input:    []byte{0x45},
			expected: 0xFFBA,
		},
		{
			name:     "two bytes",
			input:    []byte{0x45, 0x00},
			expected: 0xFFBA,
		},
		{
			name:     "ICMP header example",
			input:    []byte{0x08, 0x00, 0x00, 0x00, 0x12, 0x34, 0x00, 0x01},
			expected: 0xCAE5, // expected checksum for this header
		},
		{
			name:     "odd length",
			input:    []byte{0x45, 0x00, 0x1C},
			expected: 0xFF9E,
		},
		{
			name:     "all zeros",
			input:    []byte{0x00, 0x00, 0x00, 0x00},
			expected: 0xFFFF,
		},
		{
			name:     "all ones",
			input:    []byte{0xFF, 0xFF, 0xFF, 0xFF},
			expected: 0x0000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checksum(tt.input)
			if got != tt.expected {
				t.Errorf("checksum() = 0x%04X, want 0x%04X", got, tt.expected)
			}
		})
	}
}

// TestIcmpRunner_RunTest_InputValidation tests input validation edge cases
func TestIcmpRunner_RunTest_InputValidation(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("ICMP tests are skipped on Windows")
	}

	runner := icmpRunner{
		logger: &MockLogger{},
	}

	tests := []struct {
		name string
		sock socket.Socket
	}{
		{
			name: "whitespace only host",
			sock: socket.Socket{
				ID:   "whitespace_host",
				Name: "Whitespace Host",
				Host: "   \t\n   ",
			},
		},
		{
			name: "host with special characters",
			sock: socket.Socket{
				ID:   "special_chars_host",
				Name: "Special Characters Host",
				Host: "test@#$%^&*()host.com",
			},
		},
		{
			name: "extremely long hostname",
			sock: socket.Socket{
				ID:   "long_hostname",
				Name: "Long Hostname",
				Host: "a" + string(make([]byte, 300)) + ".com",
			},
		},
		{
			name: "hostname with unicode",
			sock: socket.Socket{
				ID:   "unicode_hostname",
				Name: "Unicode Hostname",
				Host: "тест.рф", // cyrillic domain
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			got := runner.RunTest(ctx, tt.sock)

			if got.Passed {
				t.Errorf("expected test to fail for invalid input %s, but it passed", tt.name)
			}

			if got.Error == nil {
				t.Errorf("expected error for invalid input %s, but got nil", tt.name)
			}
		})
	}
}

// TestIcmpRunner_RunTest_IPv4AddressFormats tests various IPv4 address formats
func TestIcmpRunner_RunTest_IPv4AddressFormats(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("ICMP tests are skipped on Windows")
	}

	runner := icmpRunner{
		logger: &MockLogger{},
	}

	tests := []struct {
		name       string
		host       string
		shouldPass bool
	}{
		{
			name:       "standard IPv4",
			host:       "127.0.0.1",
			shouldPass: true,
		},
		{
			name:       "IPv4 with leading zeros",
			host:       "127.000.000.001",
			shouldPass: true, // Should resolve to 127.0.0.1
		},
		{
			name:       "invalid IPv4 - too many octets",
			host:       "127.0.0.1.1",
			shouldPass: false,
		},
		{
			name:       "invalid IPv4 - octet out of range",
			host:       "256.0.0.1",
			shouldPass: false,
		},
		{
			name:       "invalid IPv4 - negative octet",
			host:       "127.0.0.-1",
			shouldPass: false,
		},
		{
			name:       "invalid IPv4 - non-numeric",
			host:       "127.0.0.x",
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sock := socket.Socket{
				ID:   tt.name,
				Name: tt.name,
				Host: tt.host,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			got := runner.RunTest(ctx, sock)

			if tt.shouldPass && !got.Passed {
				t.Logf("expected pass but got failure for %s: %v", tt.name, got.Error)
			} else if !tt.shouldPass && got.Passed {
				t.Errorf("expected failure but got pass for %s", tt.name)
			}
		})
	}
}

// TestIcmpRunner_RunTest_DNSResolutionEdgeCases tests DNS resolution edge cases
func TestIcmpRunner_RunTest_DNSResolutionEdgeCases(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("ICMP tests are skipped on Windows")
	}

	runner := icmpRunner{
		logger: &MockLogger{},
	}

	sock := socket.Socket{
		ID:   "ipv6_capable_domain",
		Name: "IPv6 Capable Domain",
		Host: "ipv6.google.com",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	got := runner.RunTest(ctx, sock)

	if !got.Passed {
		t.Logf("IPv6 capable domain failed (expected if no IPv4): %v", got.Error)
	}
}

// TestRunSocketTest is an integration test. It executes network calls to
// external public servers.
func TestRunSocketTest(t *testing.T) {
	t.Run("output chan is closed and the wait group is not blocking after a successful concurrent test", func(t *testing.T) {
		sock := socket.Socket{
			ID:   "google_tcp",
			Name: "Google TCP",
			Host: "google.com",
			Port: 80,
		}

		want := socket.Result{
			Socket: sock,
			Passed: true,
		}

		c := make(chan socket.Result)
		wg := &sync.WaitGroup{}
		cfg, err := config.NewConfig(flag.CommandLine, []string{"--timeout=1", "--verbose=false", "mocksource.json"})
		if err != nil {
			t.Fatalf("unexpected error creating config: %v", err)
		}
		done := make(chan struct{})

		wg.Add(1)
		go RunSocketTest(sock, c, wg, cfg, &MockLogger{})

		go func() {
			wg.Wait()
			done <- struct{}{}
		}()

		got := <-c

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatalf("RunSocketTest: timed out waiting for the test results")
		}

		select {
		// Once the test is finished no further results are sent.
		// If this select case blocks instead of reading the default value immediately then the channel is not closed.
		case <-c:
		default:
			t.Error("RunSocketTest: the output channel has not been closed after returning")
		}

		if !cmp.Equal(got, want) {
			t.Fatalf("RunSocketTest:\n want = %v\n got = %v\n", want, got)
		}
	})
}

func TestNewNetRunner(t *testing.T) {
	type args struct {
		sock   socket.Socket
		logger logger.Logger
	}

	tests := []struct {
		name    string
		args    args
		want    NetRunner
		wantErr bool
	}{
		{
			name: "returns an error on an empty socket",
			args: args{
				sock:   socket.Socket{},
				logger: &MockLogger{},
			},
			wantErr: true,
		},
		{
			name: "returns an httpRunner when given an HTTPs socket",
			args: args{
				sock: socket.Socket{
					ID:                "google_https",
					Name:              "Google HTTPs",
					Host:              "https://google.com",
					Port:              443,
					ExpectedHTTPCodes: []int{200},
					PathHTTP:          "/",
				},
				logger: &MockLogger{},
			},
			want: &httpRunner{
				client: &http.Client{},
				logger: &MockLogger{},
			},
			wantErr: false,
		},
		{
			name: "returns an httpRunner when given a HTTP socket",
			args: args{
				sock: socket.Socket{
					ID:                "google_http",
					Name:              "Google HTTP",
					Host:              "http://www.google.com",
					Port:              80,
					ExpectedHTTPCodes: []int{200},
					PathHTTP:          "/",
				},
				logger: &MockLogger{},
			},
			want: &httpRunner{
				client: &http.Client{},
				logger: &MockLogger{},
			},
			wantErr: false,
		},
		{
			name: "returns a tcpRunner when given a TCP socket",
			args: args{
				sock: socket.Socket{
					Port:              80,
					ExpectedHTTPCodes: []int{200},
					PathHTTP:          "/",
				},
				logger: &MockLogger{},
			},
			want: &tcpRunner{
				logger: &MockLogger{},
			},
			wantErr: false,
		},
		{
			name: "returns an icmpRunner when given an ICMP socket",
			args: args{
				sock: socket.Socket{
					Host: "google.com",
				},
				logger: &MockLogger{},
			},
			want: &icmpRunner{
				logger: &MockLogger{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNetRunner(tt.args.sock, &MockLogger{})
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewNetRunner():\n error = %v\n wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("NewNetRunner():\n got = %v\n want = %v", got, tt.want)
			}
		})
	}
}

// TestTcpRunner_RunTest is an integration test. It executes network calls to
// external public servers.
func TestTcpRunner_RunTest(t *testing.T) {
	type fields struct {
		verbose bool
	}
	type args struct {
		sock socket.Socket
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   socket.Result
	}{
		{
			name: "returns a success on a call to a valid TCP server",
			fields: fields{
				verbose: testing.Verbose(),
			},
			args: args{
				sock: socket.Socket{
					ID:   "google_tcp",
					Name: "Google TCP",
					Host: "google.com",
					Port: 80,
				},
			},
			want: socket.Result{
				Socket: socket.Socket{
					ID:   "google_tcp",
					Name: "Google TCP",
					Host: "google.com",
					Port: 80,
				},
				Passed: true,
			},
		},
		{
			name: "returns an error when the TCP connection fails",
			fields: fields{
				verbose: testing.Verbose(),
			},
			args: args{
				sock: socket.Socket{
					ID:   "invalid_tcp",
					Name: "Invalid TCP endpoint",
					Host: "doesnotexist.invalid",
					Port: 80,
				},
			},
			want: socket.Result{
				Socket: socket.Socket{
					ID:   "invalid_tcp",
					Name: "Invalid TCP endpoint",
					Host: "doesnotexist.invalid",
					Port: 80,
				},
				Passed: false,
				Error:  cmpopts.AnyError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tcpRunner{
				&MockLogger{},
			}

			if got := r.RunTest(context.Background(), tt.args.sock); !cmp.Equal(got, tt.want, cmpopts.EquateErrors()) {
				t.Fatalf("tcpRunner.RunTest():\n got = %v\n want = %v", got, tt.want)
			}
		})
	}
}

// TestHttpRunner_RunTest is an integration test. It executes network calls to
// external public servers.
func TestHttpRunner_RunTest(t *testing.T) {
	type args struct {
		sock socket.Socket
	}
	tests := []struct {
		name   string
		runner httpRunner
		args   args
		want   socket.Result
	}{
		{
			name: "returns a success on a call to a valid HTTPs server",
			runner: httpRunner{
				client: &http.Client{},
				logger: &MockLogger{},
			},
			args: args{
				sock: socket.Socket{
					ID:                "google_http",
					Name:              "Google HTTP",
					Host:              "https://www.google.com",
					Port:              443,
					ExpectedHTTPCodes: []int{200},
					PathHTTP:          "/",
				},
			},
			want: socket.Result{
				Socket: socket.Socket{
					ID:                "google_http",
					Name:              "Google HTTP",
					Host:              "https://www.google.com",
					Port:              443,
					ExpectedHTTPCodes: []int{200},
					PathHTTP:          "/",
				},
				Passed:       true,
				ResponseCode: 200,
			},
		},
		{
			name: "returns a failure on a call to an invalid HTTPs server",
			// Since both DNS and HTTPs use TCP, the conn opens successfully but,
			// the request timeouts while awaiting HTTP headers.
			runner: httpRunner{
				client: &http.Client{Timeout: time.Second},
				logger: &MockLogger{},
			},
			args: args{
				sock: socket.Socket{
					ID:                "cloudflare_dns",
					Name:              "Cloudflare DNS",
					Host:              "https://1.1.1.1",
					Port:              53,
					ExpectedHTTPCodes: []int{200},
					PathHTTP:          "/",
				},
			},
			want: socket.Result{
				Socket: socket.Socket{
					ID:                "cloudflare_dns",
					Name:              "Cloudflare DNS",
					Host:              "https://1.1.1.1",
					Port:              53,
					ExpectedHTTPCodes: []int{200},
					PathHTTP:          "/",
				},
				Passed: false,
				Error:  cmpopts.AnyError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.runner.RunTest(context.Background(), tt.args.sock)
			if !cmp.Equal(got, tt.want, cmpopts.EquateErrors()) {
				t.Fatalf("httpRunner.RunTest():\n got = %v\n want = %v", got, tt.want)
			}
		})
	}
}

// TestIcmpRunner_RunTest is an integration test. It executes network calls to
// external public servers.
// This test is common for all OS implementations except for Windows which is not supported.
func TestIcmpRunner_RunTest(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("ICMP tests are skipped on Windows")
	}

	type args struct {
		sock socket.Socket
	}

	tests := []struct {
		name   string
		runner icmpRunner
		args   args
		want   socket.Result

		// A slice defining OSes on which the given test should be skipped
		skipOn []string
	}{
		{
			name: "returns a success on a call to a valid host",
			runner: icmpRunner{
				logger: &MockLogger{},
			},
			args: args{
				sock: socket.Socket{
					ID:   "google_icmp",
					Name: "Google ICMP",
					Host: "google.com",
				},
			},
			want: socket.Result{
				Socket: socket.Socket{
					ID:   "google_icmp",
					Name: "Google ICMP",
					Host: "google.com",
				},
				Passed: true,
			},
			skipOn: []string{"darwin"},
		},
		{
			name: "returns a success on a call to a valid IP address",
			runner: icmpRunner{
				logger: &MockLogger{},
			},
			args: args{
				sock: socket.Socket{
					ID:   "google_icmp",
					Name: "Google ICMP",
					Host: "8.8.8.8",
				},
			},
			want: socket.Result{
				Socket: socket.Socket{
					ID:   "google_icmp",
					Name: "Google ICMP",
					Host: "8.8.8.8",
				},
				Passed: true,
			},
			skipOn: []string{"darwin"},
		},
		{
			name: "returns an error on an empty host",
			runner: icmpRunner{
				logger: &MockLogger{},
			},
			args: args{
				sock: socket.Socket{
					ID:   "empty_host",
					Name: "Empty Host",
					Host: "",
				},
			},
			want: socket.Result{
				Socket: socket.Socket{
					ID:   "empty_host",
					Name: "Empty Host",
					Host: "",
				},
				Passed: false,
				Error:  cmpopts.AnyError,
			},
		},
		{
			name: "returns an error on an invalid IP address",
			runner: icmpRunner{
				logger: &MockLogger{},
			},
			args: args{
				sock: socket.Socket{
					ID:   "invalid_ip",
					Name: "Invalid IP",
					Host: "256.100.50.25",
				},
			},
			want: socket.Result{
				Socket: socket.Socket{
					ID:   "invalid_ip",
					Name: "Invalid IP",
					Host: "256.100.50.25",
				},
				Passed: false,
				Error:  cmpopts.AnyError,
			},
		},
		// {
		// 	name: "returns a success on a call to localhost using hostname",
		// 	runner: icmpRunner{
		// 		logger: &MockLogger{},
		// 	},
		// 	args: args{
		// 		sock: socket.Socket{
		// 			ID:   "localhost_icmp",
		// 			Name: "Localhost ICMP",
		// 			Host: "localhost",
		// 		},
		// 	},
		// 	want: socket.Result{
		// 		Socket: socket.Socket{
		// 			ID:   "localhost_icmp",
		// 			Name: "Localhost ICMP",
		// 			Host: "localhost",
		// 		},
		// 		Passed: true,
		// 	},
		// },
		{
			name: "returns a success on a call to localhost using IP",
			runner: icmpRunner{
				logger: &MockLogger{},
			},
			args: args{
				sock: socket.Socket{
					ID:   "localhost_icmp",
					Name: "Localhost ICMP",
					Host: "127.0.0.1",
				},
			},
			want: socket.Result{
				Socket: socket.Socket{
					ID:   "localhost_icmp",
					Name: "Localhost ICMP",
					Host: "127.0.0.1",
				},
				Passed: true,
			},
		},
	}
	for _, tt := range tests {
		if tt.skipOn != nil && slices.Contains(tt.skipOn, runtime.GOOS) {
			t.Logf("skipping test %s on %s", tt.name, runtime.GOOS)
			continue
		}

		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			got := tt.runner.RunTest(ctx, tt.args.sock)
			if !cmp.Equal(got, tt.want, cmpopts.EquateErrors()) {
				t.Fatalf("icmpRunner.RunTest():\n got = %v\n want = %v", got, tt.want)
			}
		})
	}
}

func TestIcmpRunner_RunTest_Windows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip()
	}

	socket := socket.Socket{
		ID:   "google_icmp",
		Name: "Google ICMP",
		Host: "google.com",
	}

	runner, err := NewNetRunner(socket, &MockLogger{})
	if err != nil {
		t.Error("failed to create a new Windows ICMP runner")
	}

	result := runner.RunTest(context.Background(), socket)
	if result.Error == nil {
		t.Error("expected error, got nil")
	}
}

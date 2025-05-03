package netrunner

import (
	"context"
	"flag"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/socket"
)

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
		go RunSocketTest(sock, c, wg, cfg)

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
		verbose bool
		sock    socket.Socket
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
				verbose: testing.Verbose(),
				sock:    socket.Socket{},
			},
			wantErr: true,
		},
		{
			name: "returns an httpRunner when given an HTTPs socket",
			args: args{
				verbose: testing.Verbose(),
				sock: socket.Socket{
					ID:                "google_https",
					Name:              "Google HTTPs",
					Host:              "https://google.com",
					Port:              443,
					ExpectedHTTPCodes: []int{200},
					PathHTTP:          "/",
				},
			},
			want: httpRunner{
				client:  &http.Client{},
				verbose: testing.Verbose(),
			},
			wantErr: false,
		},
		{
			name: "returns an httpRunner when given a HTTP socket",
			args: args{
				verbose: testing.Verbose(),
				sock: socket.Socket{
					ID:                "google_http",
					Name:              "Google HTTP",
					Host:              "http://www.google.com",
					Port:              80,
					ExpectedHTTPCodes: []int{200},
					PathHTTP:          "/",
				},
			},
			want: httpRunner{
				client:  &http.Client{},
				verbose: testing.Verbose(),
			},
			wantErr: false,
		},
		{
			name: "returns a tcpRunner when given a TCP socket",
			args: args{
				verbose: testing.Verbose(),
				sock: socket.Socket{
					Port:              80,
					ExpectedHTTPCodes: []int{200},
					PathHTTP:          "/",
				},
			},
			want:    tcpRunner{verbose: testing.Verbose()},
			wantErr: false,
		},
		{
			name: "returns an icmpRunner when given an ICMP socket",
			args: args{
				verbose: testing.Verbose(),
				sock: socket.Socket{
					Host: "google.com",
				},
			},
			want:    icmpRunner{verbose: testing.Verbose()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNetRunner(tt.args.sock, tt.args.verbose)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tcpRunner{tt.fields.verbose}

			if got := r.RunTest(context.Background(), tt.args.sock); !cmp.Equal(got, tt.want) {
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
			name:   "returns a success on a call to a valid HTTPs server",
			runner: httpRunner{client: &http.Client{}, verbose: testing.Verbose()},
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
			runner: httpRunner{client: &http.Client{Timeout: time.Second}, verbose: testing.Verbose()},
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
// This test is common for all OS implementations.
func TestIcmpRunner_RunTest(t *testing.T) {
	type args struct {
		sock socket.Socket
	}
	tests := []struct {
		name   string
		runner icmpRunner
		args   args
		want   socket.Result
	}{
		{
			name: "returns a success on a call to a valid host",
			runner: icmpRunner{
				verbose: testing.Verbose(),
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
		},
		{
			name: "returns a success on a call to a valid IP address",
			runner: icmpRunner{
				verbose: testing.Verbose(),
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
		},
		{
			name: "returns an error on an empty host",
			runner: icmpRunner{
				verbose: testing.Verbose(),
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
				verbose: testing.Verbose(),
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
	}
	for _, tt := range tests {
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

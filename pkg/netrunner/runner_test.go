package netrunner

import (
	"context"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
		done := make(chan struct{})

		wg.Add(1)
		go RunSocketTest(sock, c, wg, 1, false)

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
			name: "returns a httpRunner when given an HTTPs socket",
			args: args{
				verbose: false,
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
				verbose: false,
			},
			wantErr: false,
		},
		{
			name: "returns a httpRunner when given a HTTP socket",
			args: args{
				verbose: false,
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
				verbose: false,
			},
			wantErr: false,
		},
		{
			name: "returns a tcpRunner when given a TCP socket",
			args: args{
				verbose: false,
				sock: socket.Socket{
					ID:                "",
					Name:              "",
					Host:              "",
					Port:              80,
					ExpectedHTTPCodes: []int{200},
					PathHTTP:          "/",
				},
			},
			want:    tcpRunner{verbose: false},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNetRunner(tt.args.sock, tt.args.verbose)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewNetRunner():\n error = %v\n wantErr = %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNetRunner():\n got = %v\n want = %v", got, tt.want)
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
				verbose: false,
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
				t.Errorf("tcpRunner.RunTest():\n got = %v\n want = %v", got, tt.want)
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
			runner: httpRunner{client: &http.Client{}, verbose: false},
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
			// The since both DNS and HTTPs use TCP the conn opens successfully but
			// the request timeouts while awaiting HTTP headers.
			runner: httpRunner{client: &http.Client{Timeout: time.Second}, verbose: false},
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
				t.Errorf("httpRunner.RunTest():\n got = %v\n want = %v", got, tt.want)
			}
		})
	}
}

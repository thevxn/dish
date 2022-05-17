package socket

import ( 
	"testing"
)

func FetchSocketList(t *testing.T) {
	var want Sockets
	if got := GibPole("demo_sockets.json", false); got != want {
		t.Errorf("GibPole() = %q, want %q", got, want)
	}
}

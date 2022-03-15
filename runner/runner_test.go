package runner

import (
	"testing"
)

func TestCheckSite(t *testing.T) {
        want := 200
        if got := CheckSite("http://traefik.savla.su", 80); got != want {
                t.Errorf("GibPole() = %q, want %q", got, want)
        }
}


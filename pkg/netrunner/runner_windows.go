//go:build windows

package netrunner

import (
	"context"
	"errors"

	"go.vxn.dev/dish/pkg/socket"
)

type icmpRunner struct {
	verbose bool
}

func (runner icmpRunner) RunTest(ctx context.Context, sock socket.Socket) socket.Result {
	return socket.Result{Socket: sock, Error: errors.New("icmp tests on windows are not implemented")}
}

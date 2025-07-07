//go:build windows

package netrunner

import (
	"context"
	"errors"

	"go.vxn.dev/dish/pkg/logger"
	"go.vxn.dev/dish/pkg/socket"
)

type icmpRunner struct {
	logger logger.Logger
}

func (runner icmpRunner) RunTest(ctx context.Context, sock socket.Socket) socket.Result {
	return socket.Result{Socket: sock, Error: errors.New("icmp tests on windows are not implemented")}
}

func checksum(data []byte) uint16 {
    return 0 // return invalid checksum since not implemented in Windows
}

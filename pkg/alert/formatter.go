package alert

import (
	"fmt"

	"go.vxn.dev/dish/pkg/socket"
)

func FormatMessengerText(result socket.Result) string {
	// Hotfix unsupported <nil> tag by TG
	if result.Error == nil {
		result.Error = fmt.Errorf("")
	}

	if result.Socket.PathHTTP != "" {
		return fmt.Sprintf("• %s:%d%s -- %v\n",
			result.Socket.Host, result.Socket.Port, result.Socket.PathHTTP, result.Error)
	}
	return fmt.Sprintf("• %s:%d -- %v\n", result.Socket.Host, result.Socket.Port, result.Error)
}

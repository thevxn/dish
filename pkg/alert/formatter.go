package alert

import (
	"fmt"

	"go.vxn.dev/dish/pkg/socket"
)

func FormatMessengerText(result socket.Result) string {

	status := "failed"
	if result.Passed {
		status = "success"
	}

	text := fmt.Sprintf("• %s:%d", result.Socket.Host, result.Socket.Port)

	if result.Socket.PathHTTP != "" {
		text += result.Socket.PathHTTP
	}

	text += " -- " + status

	if status == "failed" {
		text += " -- "
		text += result.Error.Error()
	}

	text += "\n"

	return text
}

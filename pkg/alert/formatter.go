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
		text += " \u274C" // ❌
		text += " -- "
		text += result.Error.Error()
	} else {
		text += " \u2705" // ✅
	}

	text += "\n"

	return text
}

func FormatMessengerTextWithHeader(header, body string) string {
	return header + "\n\n" + body
}

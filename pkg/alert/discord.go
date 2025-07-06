package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/logger"
)

type discordSender struct {
	botToken      string
	channelID     string
	httpClient    HTTPClient
	logger        logger.Logger
	notifySuccess bool
}

type discordMessagePayload struct {
	Content string `json:"content"`
}

const (
	discordBaseURL      = "https://discord.com/api/v10"
	discordMessagesPath = "/channels/%s/messages"
	discordMessagesURL  = discordBaseURL + discordMessagesPath
)

func NewDiscordSender(httpClient HTTPClient, config *config.Config, logger logger.Logger) ChatNotifier {
	return &discordSender{
		botToken:      config.DiscordBotToken,
		channelID:     config.DiscordChannelID,
		httpClient:    httpClient,
		logger:        logger,
		notifySuccess: config.TextNotifySuccess,
	}

}

func (s *discordSender) send(message string, failedCount int) error {
	// If no checks failed and success should not be notified, there is nothing to send
	if failedCount == 0 && !s.notifySuccess {
		s.logger.Debug("no sockets failed, nothing will be sent to Telegram")

		return nil
	}

	formattedSendURL := fmt.Sprintf(discordMessagesURL, s.channelID)

	payload := discordMessagePayload{Content: message}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error submitting discord alert: %w ", err)
	}

	resp, err := handleSubmit(s.httpClient, http.MethodPost, formattedSendURL, bytes.NewBuffer(body), func(o *submitOptions) {
		o.headers["Authorization"] = "Bot " + strings.TrimSpace(s.botToken)
	})

	if err != nil {
		return fmt.Errorf("error submitting discord alert: %w", err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("error submitting discord alert: non-success status code: %d", resp.StatusCode)
	}

	s.logger.Debug("discord message sent")
	return nil
}

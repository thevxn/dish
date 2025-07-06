package alert

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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

const (
	discordBaseURL         = "https://discord.com/api/v10"
	discordSendMessagePath = "/channels/%s/messages"
	discordSendMessageURL  = discordBaseURL + discordSendMessagePath
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

	formattedSendURL := fmt.Sprintf(discordSendMessageURL, s.channelID)

	payload := map[string]string{
		"content": message,
	}
	body, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("error submitting discord alert: %w ", err)
	}

	req, err := http.NewRequest("POST", formattedSendURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error submitting discord alert: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", s.botToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error submitting discord alert: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return errors.New("error submitting discord alert: non-success status code")
	}

	s.logger.Debug("discord message sent")
	return nil
}

// Test token : AR7DbyW7CL3HxQiDeyLq-aZfaL9jnV_8
// Test channel: 1391401731329622036/1391401732046585888

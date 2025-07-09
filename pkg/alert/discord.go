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

const discordMessageTitle = "📡 **dish run results**:"

type discordSender struct {
	botToken      string
	channelID     string
	httpClient    HTTPClient
	logger        logger.Logger
	notifySuccess bool
	url           string
}

type discordMessagePayload struct {
	Content string `json:"content"`
	Flags   int    `json:"flags,omitempty"`
}

const (
	discordBaseURL      = "https://discord.com/api/v10"
	discordMessagesPath = "/channels/%s/messages"
	discordMessagesURL  = discordBaseURL + discordMessagesPath
)

func NewDiscordSender(
	httpClient HTTPClient,
	config *config.Config,
	logger logger.Logger,
) (ChatNotifier, error) {
	parsedURL, err := parseAndValidateURL(
		fmt.Sprintf(discordMessagesURL, strings.TrimSpace(config.DiscordChannelID)),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &discordSender{
		botToken:      config.DiscordBotToken,
		channelID:     config.DiscordChannelID,
		httpClient:    httpClient,
		logger:        logger,
		notifySuccess: config.TextNotifySuccess,
		url:           parsedURL.String(),
	}, nil
}

func (s *discordSender) send(message string, failedCount int) error {
	// If no checks failed and success should not be notified, there is nothing to send
	if failedCount == 0 && !s.notifySuccess {
		s.logger.Debug("no sockets failed, nothing will be sent to Telegram")

		return nil
	}

	payload := discordMessagePayload{
		Content: FormatMessengerTextWithHeader(discordMessageTitle, message),
		Flags:   4, // Suppress embedded links in the message
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error submitting discord alert: %w ", err)
	}

	resp, err := handleSubmit(
		s.httpClient,
		http.MethodPost,
		s.url,
		bytes.NewBuffer(body),
		func(o *submitOptions) {
			o.headers["Authorization"] = "Bot " + strings.TrimSpace(s.botToken)
		},
	)
	if err != nil {
		return fmt.Errorf("error submitting discord alert: %w", err)
	}

	err = handleRead(resp, s.logger)
	if err != nil {
		return fmt.Errorf(
			"error submitting discord alert: non-success status code: %d",
			resp.StatusCode,
		)
	}

	s.logger.Debug("discord message sent")
	return nil
}

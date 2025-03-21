package alert

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/message"
)

type ChatNotifier interface {
	Send(string, int) error
}
type MachineNotifier interface {
	Send(message.Results, int) error
}

type notifier struct {
	notificationSenders []ChatNotifier
	payloadSenders      []MachineNotifier
}

type telegramSender struct {
	httpClient *http.Client
	failedOnly bool
}

// TODO: failedOnly should be moved to config, hardcoded for now
func NewTelegramSender(h *http.Client) *telegramSender {
	return &telegramSender{httpClient: h, failedOnly: true}
}

func (s *telegramSender) Send(rawMessage string, failedCount int) error {
	if rawMessage == "" {
		return errors.New("empty message string given")
	}

	// If there are no failed sockets and we only wish to be notified when they fail, there is nothing to do
	if failedCount == 0 && s.failedOnly {
		log.Printf("%T: no failed sockets and failedOnly == true, nothing will be sent", s)
		return nil
	}

	// form the Telegram URL
	telegramURL := "https://api.telegram.org/bot" + config.TelegramBotToken + "/sendMessage?chat_id=" + config.TelegramChatID + "&disable_web_page_preview=True&parse_mode=HTML&text="

	msg := "<b>dish run results</b>:\n\n" + rawMessage

	// escape dish report string for Telegram
	msg = url.QueryEscape(msg)

	resp, err := s.httpClient.Get(telegramURL + msg)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// write to console log if verbose flag set
	if config.Verbose {
		log.Println(telegramURL)
		log.Println(string(body))
	}

	return nil
}

type webhookSender struct {
	httpClient *http.Client
	failedOnly bool
}

// TODO: failedOnly should be moved to config, hardcoded for now
func NewWebhookSender(h *http.Client) *webhookSender {
	return &webhookSender{httpClient: h, failedOnly: true}
}

func (s *webhookSender) Send(m message.Results, failedCount int) error {
	// If there are no failed sockets and we only wish to be notified when they fail, there is nothing to do
	if failedCount == 0 && s.failedOnly {
		log.Printf("%T: no failed sockets and failedOnly == true, nothing will be sent", s)
		return nil
	}

	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	if config.Verbose {
		log.Printf("Prepared webhook data: %v", string(jsonData))
	}

	resp, err := s.httpClient.Post(config.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if config.Verbose {
		log.Printf("Webhook notification sent. Webhook URL: %s", config.WebhookURL)
		log.Printf("Received response from webhook URL. Status: %s. Body: %s", resp.Status, string(body))
	}

	return nil
}

type apiSender struct {
	httpClient *http.Client
	failedOnly bool
}

// TODO: failedOnly should be moved to config, hardcoded for now
func NewApiSender(h *http.Client) *apiSender {
	return &apiSender{httpClient: h, failedOnly: false}
}

func (s *apiSender) Send(m message.Results, failedCount int) error {
	// If there are no failed sockets and we only wish to be notified when they fail, there is nothing to do
	if failedCount == 0 && s.failedOnly {
		log.Printf("%T: no failed sockets and failedOnly == true, nothing will be sent", s)
		return nil
	}

	jsonData, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	bodyReader := bytes.NewReader(jsonData)
	log.Println(string(jsonData))

	url := config.UpdateURL

	regex, err := regexp.Compile("^(http|https)://")
	if err != nil {
		return err
	}
	match := regex.MatchString(url)

	if !match {
		return nil
	}

	// Push results
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(config.HeaderName, config.HeaderValue)

	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	log.Println("Results pushed to swapi")

	return nil
}

// TODO: Inject config after it has been refactored
func NewNotifier(httpClient *http.Client) *notifier {
	// Set chat integrations to be notified (e.g. Telegram)
	notificationSenders := make([]ChatNotifier, 0)
	if config.UseTelegram {
		notificationSenders = append(notificationSenders, NewTelegramSender(httpClient))
	}

	// Set machine interface integrations to be notified (e.g. Webhooks)
	payloadSenders := make([]MachineNotifier, 0)
	if config.UpdateStates {
		payloadSenders = append(payloadSenders, NewApiSender(httpClient))
	}
	if config.UseWebhooks {
		payloadSenders = append(payloadSenders, NewWebhookSender(httpClient))
	}

	return &notifier{notificationSenders, payloadSenders}
}

func (n *notifier) SendChatNotifications(m string, failedCount int) error {
	var errs []error

	if len(n.notificationSenders) == 0 {
		if config.Verbose {
			log.Println("no chat notification receivers configured, no notifications sent.")
		}
		return nil
	}

	for _, sender := range n.notificationSenders {
		if err := sender.Send(m, failedCount); err != nil {
			log.Printf("failed to send notification using %T: %v", sender, err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to send chat notifications: %w", errors.Join(errs...))
	}

	return nil
}

func (n *notifier) SendMachineNotifications(m message.Results, failedCount int) error {
	var errs []error

	if len(n.payloadSenders) == 0 {
		if config.Verbose {
			log.Println("no machine interface payload receivers configured, no notifications sent.")
		}
		return nil
	}
	for _, sender := range n.payloadSenders {
		if err := sender.Send(m, failedCount); err != nil {
			log.Printf("failed to send notification using %T: %v", sender, err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to send machine notifications: %w", errors.Join(errs...))
	}

	return nil
}

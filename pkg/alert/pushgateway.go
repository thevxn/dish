package alert

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	jobName         = "dish_results"
	failedCountName = "dish_failed_count"
	failedCountHelp = "#HELP failed sockets registered by dish"
	failedCountType = "#TYPE dish_failed_count counter"
)

type pushgatewaySender struct {
	httpClient    *http.Client
	url           string
	instanceName  string
	verbose       bool
	notifySuccess bool
}

func NewPushgatewaySender(httpClient *http.Client, url string, instanceName string, verbose bool, notifySuccess bool) *pushgatewaySender {
	return &pushgatewaySender{
		httpClient,
		url,
		instanceName,
		verbose,
		notifySuccess,
	}
}

// createMessage returns a string containing the message text in Pushgateway-specific format.
func (s *pushgatewaySender) createMessage(failedCount int) string {
	msg := fmt.Sprintln(failedCountHelp)
	msg += fmt.Sprintln(failedCountType)
	msg += fmt.Sprintln(failedCountName, strconv.Itoa(failedCount))

	return msg
}

// Send pushes the results to Pushgateway.
//
// The first argument is needed to implement the MachineNotifier interface, however, it is ignored in favor of a custom message implementation via the createMessage method.
func (s *pushgatewaySender) send(_ Results, failedCount int) error {
	// If no checks failed and failedOnly is set to true, there is nothing to send
	if failedCount == 0 && !s.notifySuccess {
		if s.verbose {
			log.Println("no sockets failed and notifySuccess == false, nothing will be sent to Pushgateway")
		}
		return nil
	}

	msg := s.createMessage(failedCount)

	bodyReader := bytes.NewReader([]byte(msg))

	formattedURL := s.url + "/metrics/job/" + jobName + "/instance/" + s.instanceName

	req, err := http.NewRequest(http.MethodPut, formattedURL, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/byte")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response code received from Pushgateway (expected: %d, got: %d)", http.StatusOK, res.StatusCode)
	}

	// Write the body to console if verbose flag set
	if s.verbose {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
		log.Println("pushgateway response:", string(body))
	}

	log.Println("results pushed to Pushgateway")

	return nil
}

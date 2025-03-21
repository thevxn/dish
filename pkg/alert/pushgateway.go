package alert

import (
	"bytes"
	"fmt"
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
	httpClient   *http.Client
	url          string
	instanceName string
}

func NewPushgatewaySender(httpClient *http.Client, url string, instanceName string) *pushgatewaySender {
	return &pushgatewaySender{
		httpClient,
		url,
		instanceName,
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

	msg := s.createMessage(failedCount)

	bodyReader := bytes.NewReader([]byte(msg))

	formattedURL := s.url + "/metrics/job/" + jobName + "/instance/" + s.instanceName

	req, err := http.NewRequest(http.MethodPut, formattedURL, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/byte")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	log.Println("Results pushed to pushgateway")

	return nil
}

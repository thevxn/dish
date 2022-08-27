package reporter

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"savla-dish/pkg/config"
	"strconv"
)

const (
	jobName         = "dish_results"
	instanceName    = "generic-dish"
	failedCountName = "dish_failed_count"
	failedCountHelp = "#HELP failed sockets registered by savla-dish"
	failedCountType = "#TYPE dish_failed_count counter"
)

type Message struct {
	FailedCount int
	body        string
}

func MakeMessage(count int) Message {
	msg := Message{}
	msg.FailedCount = count
	messageString := fmt.Sprintln(failedCountName, strconv.Itoa(msg.FailedCount))
	msg.body = messageString
	return msg
}

func (msg Message) PushDishResults() error {
	bodyReader := bytes.NewReader([]byte(msg.body))
	formattedURL := config.TargetURL + "/metrics/job/" + jobName + "/instance/" + instanceName

	log.Println(formattedURL)

	req, err := http.NewRequest(http.MethodPost, formattedURL, bodyReader)
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

	return nil
}

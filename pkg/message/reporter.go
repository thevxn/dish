package message

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"savla-dish/pkg/config"
)

const (
	jobName         = "dish_results"
	failedCountName = "dish_failed_count"
	failedCountHelp = "#HELP failed sockets registered by savla-dish"
	failedCountType = "#TYPE dish_failed_count counter"
)

type Message struct {
	FailedCount int
	body        string
}

func Make(count int) Message {
	msg := Message{}
	msg.FailedCount = count

	// ensure HELP and TYPE fields are added too!
	messageString := fmt.Sprintln(failedCountHelp)
	messageString += fmt.Sprintln(failedCountType)
	messageString += fmt.Sprintln(failedCountName, strconv.Itoa(msg.FailedCount))

	msg.body = messageString
	return msg
}

func (msg Message) PushDishResults() error {
	bodyReader := bytes.NewReader([]byte(msg.body))
	formattedURL := config.TargetURL + "/metrics/job/" + jobName + "/instance/" + config.InstanceName

	log.Println(formattedURL)

	// push requests use PUT method
	req, err := http.NewRequest("PUT", formattedURL, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/byte")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	log.Println("Results pushed to pushgateway")
	defer res.Body.Close()

	return nil
}

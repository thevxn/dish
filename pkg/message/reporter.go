package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"dish/pkg/config"
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

type Results struct {
	Map map[string]bool `json:"dish_results"`
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

func UpdateSocketStates(results Results) error {
	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil
	}
	bodyReader := bytes.NewReader([]byte(jsonData))

	url := config.UpdateURL

	regex, err := regexp.Compile("^(http|https)://")
	match := regex.MatchString(url)

	if !match {
		return nil
	}

	// push requests use PUT method
	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/byte")
	req.Header.Set(config.HeaderName, config.HeaderValue)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	log.Println("Results pushed to swapi")
	defer res.Body.Close()

	return nil
}

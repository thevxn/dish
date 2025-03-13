package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"go.vxn.dev/dish/pkg/config"
)

const (
	jobName         = "dish_results"
	failedCountName = "dish_failed_count"
	failedCountHelp = "#HELP failed sockets registered by dish"
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

	// Push requests use PUT method
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

func UpdateSocketStates(results Results) error {
	jsonData, err := json.Marshal(results)
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

	// Push requests use PUT method
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(config.HeaderName, config.HeaderValue)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	log.Println("Results pushed to swapi")

	return nil
}

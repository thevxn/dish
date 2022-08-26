package reporter

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	jobName         = "dish_results"
	instanceName    = "generic-dish"
	failedCountName = "dish_failed_count"
	failedCountHelp = "#HELP failed sockets registered by savla-dish"
	failedCountType = "#TYPE dish_failed_count counter"
)

type Report struct {
	FailedCount int8
	message     []byte
	timestamp   time.Time
}

var (
	Reporter       Report
	TargetURL      *string
	UsePushgateway *bool
)

func composeMessage() []byte {
	var messageString string
	messageString += failedCountHelp + "\n"
	messageString += failedCountType + "\n"
	messageString += failedCountName + " " + strconv.Itoa(int(Reporter.FailedCount)) + "\n"

	log.Println(messageString)
	return []byte(messageString)
}

func PushDishResults() error {
	Reporter.message = composeMessage()

	bodyReader := bytes.NewReader(Reporter.message)
	formattedURL := *TargetURL + "/metrics/job/" + jobName + "/instance/" + instanceName

	log.Println(formattedURL)

	req, err := http.NewRequest(http.MethodPost, formattedURL, bodyReader)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/byte")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	return nil
}

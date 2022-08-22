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

func composeMessage() (message []byte) {
	var messageString string
	messageString += failedCountHelp + "\n"
	messageString += failedCountType + "\n"
	messageString += failedCountName + " " + strconv.Itoa(int(Reporter.FailedCount)) + "\n"

	log.Println(messageString)
	return []byte(messageString)
}

func PushDishResults() (err error) {
	Reporter.message = composeMessage()

	bodyReader := bytes.NewReader(Reporter.message)
	formattedURL := *TargetURL + "/metrics/job/" + jobName + "/instance/" + instanceName

	log.Println(formattedURL)

	req, err := http.NewRequest(http.MethodPost, formattedURL, bodyReader)
	if err != nil {
		panic(err)
		return err
	}
	req.Header.Set("Content-Type", "application/byte")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		panic(err)
		return err
	}

	defer res.Body.Close()

	return nil
}

// fetchFileStream
/*
func fetchFileStream(input string) (byteStream *[]byte) {
	//jsonFile, err := os.Open("sockets.json")
	jsonFile, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	defer jsonFile.Close()

	// use local var as "buffer", then return pointer to data
	stream, _ := ioutil.ReadAll(jsonFile)
	return &stream
}*/

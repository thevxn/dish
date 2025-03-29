package alert

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

type apiSender struct {
	httpClient  *http.Client
	url         string
	headerName  string
	headerValue string
}

func NewApiSender(httpClient *http.Client, url string, headerName string, headerValue string) *apiSender {
	return &apiSender{
		httpClient,
		url,
		headerName,
		headerValue,
	}
}

func (s *apiSender) send(m Results, failedCount int) error {

	jsonData, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	bodyReader := bytes.NewReader(jsonData)
	log.Println(string(jsonData))

	url := s.url

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
	req.Header.Set(s.headerName, s.headerValue)

	res, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	log.Println("Results pushed to remote api")

	return nil
}

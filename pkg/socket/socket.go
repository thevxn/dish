// Package socket provides functionality related to handling sockets, which is a structure
// representing the target endpoint/socket to be checked.
package socket

import (
	"encoding/json"
	"fmt"
	"io"

	"go.vxn.dev/dish/pkg/config"
	"go.vxn.dev/dish/pkg/logger"
)

type Result struct {
	Socket       Socket
	Passed       bool
	ResponseCode int
	Error        error
}

type SocketList struct {
	Sockets []Socket `json:"sockets"`
}

type Socket struct {
	// ID is an unique identifier of such socket.
	ID string `json:"id"`

	// Socket name, unique identificator, snake_cased.
	Name string `json:"socket_name"`

	// Remote endpoint hostname or URL.
	Host string `json:"host_name"`

	// Remote port to assemble a socket.
	Port int `json:"port_tcp"`

	// HTTP Status Codes expected when giving the endpoint a HEAD/GET request.
	ExpectedHTTPCodes []int `json:"expected_http_code_array"`

	// HTTP Path to test on Host.
	PathHTTP string `json:"path_http"`
}

// PrintSockets prints SocketList.
func PrintSockets(list *SocketList, logger logger.Logger) {
	logger.Debug("loaded sockets:")
	for _, socket := range list.Sockets {
		logger.Debugf("Host: %s, Port: %d, ExpectedHTTPCodes: %v", socket.Host, socket.Port, socket.ExpectedHTTPCodes)
	}
}

// LoadSocketList decodes a JSON encoded SocketList from the provided io.ReadCloser.
func LoadSocketList(reader io.ReadCloser) (*SocketList, error) {
	defer reader.Close()

	list := new(SocketList)
	if err := json.NewDecoder(reader).Decode(list); err != nil {
		return nil, fmt.Errorf("error decoding sockets json: %w", err)
	}

	return list, nil
}

// FetchSocketList fetches the list of sockets to be checked. 'input' should be a string like '/path/filename.json', or an HTTP URL string.
func FetchSocketList(config *config.Config, logger logger.Logger) (*SocketList, error) {
	var reader io.ReadCloser
	var err error

	fetchHandler := NewFetchHandler(logger)
	if IsFilePath(config.Source) {
		reader, err = fetchHandler.fetchSocketsFromFile(config)
	} else {
		reader, err = fetchHandler.fetchSocketsFromRemote(config)
	}

	if err != nil {
		return nil, err
	}

	return LoadSocketList(reader)
}

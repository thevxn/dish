package main

import (
	"reflect"
	"sort"
	"testing"

	"go.vxn.dev/dish/pkg/socket"
)

// compareResults is a custom comparison function to assert the results returned from the fanInChannels function are equal to the expected results
func compareResults(expected, actual []socket.Result) bool {
	sort.Slice(expected, func(i, j int) bool {
		return expected[i].ResponseCode < expected[j].ResponseCode
	})
	sort.Slice(actual, func(i, j int) bool {
		return actual[i].ResponseCode < actual[j].ResponseCode
	})

	for i := range expected {
		if !reflect.DeepEqual(expected[i].Socket, actual[i].Socket) ||
			expected[i].Passed != actual[i].Passed ||
			expected[i].ResponseCode != actual[i].ResponseCode {
			return false
		}
	}
	return true
}

func TestFanInChannels(t *testing.T) {
	testChannels := []chan socket.Result{}

	for range 3 {
		c := make(chan socket.Result)
		testChannels = append(testChannels, c)
	}

	go func() {
		for i, channel := range testChannels {
			channel <- socket.Result{
				Socket:       socket.Socket{},
				Passed:       true,
				ResponseCode: 200 + i,
			}
			close(channel)
		}
	}()

	resultingChan := fanInChannels(testChannels...)
	actual := []socket.Result{}
	for result := range resultingChan {
		actual = append(actual, result)
	}

	expected := []socket.Result{
		{
			Socket:       socket.Socket{},
			Passed:       true,
			ResponseCode: 200,
		}, {
			Socket:       socket.Socket{},
			Passed:       true,
			ResponseCode: 201,
		}, {
			Socket:       socket.Socket{},
			Passed:       true,
			ResponseCode: 202,
		},
	}

	if !compareResults(expected, actual) {
		t.Fatalf("expected: %+v, got: %+v", expected, actual)
	}
}

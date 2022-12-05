package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req)
}

func TestGetResultSuccess(t *testing.T) {

	var testcases = []struct {
		requestedURL       string
		responseStatusCode int
		responseBody       string
		expectedMD5Hash    string
		expectedURLAddress string
	}{
		{"www.google.com", http.StatusOK, "sample response", "d38cbaf4a9b7626cb2c420be3f918c11", "http://www.google.com"},
		{"www.youtube.com", http.StatusOK, "sample response from youtube", "573b3fbccf725497a4112c4842668b05", "http://www.youtube.com"},
	}

	for _, tc := range testcases {
		expectedResponse := &http.Response{
			StatusCode: tc.responseStatusCode,
			Body:       io.NopCloser(strings.NewReader(tc.responseBody)),
		}
		client := &http.Client{}
		client.Transport = RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			return expectedResponse, nil
		})

		httpClient := &HTTPClient{context.Background(), client}
		result, err := httpClient.GetResult(tc.requestedURL)
		if err != nil {
			t.Errorf("error while .. :%s", err.Error())
		}
		if result != nil {
			if (result.URLAddress) != tc.expectedURLAddress {
				t.Errorf("expected url address is '%s' got :%s", "http://www.google.com", result.URLAddress)
			}

			if (result.MD5Hash) != tc.expectedMD5Hash {
				t.Errorf("expected hash is '%s' got :%s", "d38cbaf4a9b7626cb2c420be3f918c11", result.MD5Hash)
			}
		}
	}

}

func TestGetResultError(t *testing.T) {

	var testcases = []struct {
		requestedURL       string
		responseStatusCode int
		expectedError      string
	}{
		{"www.google.com", http.StatusForbidden, "Not Allowed to access"},
		{"www.youtube.com", http.StatusInternalServerError, "Internal Server error"},
	}

	for _, tc := range testcases {
		expectedResponse := &http.Response{
			StatusCode: tc.responseStatusCode,
			Body:       io.NopCloser(strings.NewReader("error while getting the respons")),
		}
		client := &http.Client{}
		client.Transport = RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			return expectedResponse, fmt.Errorf(tc.expectedError)
		})

		httpClient := &HTTPClient{context.Background(), client}
		_, err := httpClient.GetResult(tc.requestedURL)
		if err == nil {
			t.Fatal("expected error,but got nothing")
		}
		if err != nil {
			if !strings.Contains(err.Error(), tc.expectedError) {
				t.Errorf("Expected error should contain '%s' ,but got :%s", tc.expectedError, err)
			}
		}
	}

}

func TestGetResultError_ContextCancel(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	httpClient := &HTTPClient{ctx, &http.Client{}}
	cancel()
	_, err := httpClient.GetResult("www.google.com")
	if err == nil {
		t.Fatal("expected error,but got nothing")
	}

	if err != nil {
		if !strings.Contains(err.Error(), "context canceled") {
			t.Errorf("Expected error should contain '%s' ,but got :%s", "context canceled", err)
		}
	}

}

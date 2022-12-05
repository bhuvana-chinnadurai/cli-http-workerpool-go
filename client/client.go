package client

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client interface {
	GetResult(urlStr string) (*Result, error)
}

type HTTPClient struct {
	ctx    context.Context
	client *http.Client
}

type Result struct {
	URLAddress string
	MD5Hash    string
}

func New(ctx context.Context) *HTTPClient {
	return &HTTPClient{ctx: ctx, client: &http.Client{}}
}

func (h *HTTPClient) GetResult(urlStr string) (*Result, error) {
	url := &url.URL{Scheme: "http", Host: urlStr}

	ctx, cancel := context.WithTimeout(h.ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error while preparing the request: '%s' ", err.Error())
	}

	response, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while making the request '%s' ", err.Error())
	}
	defer response.Body.Close()
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading the response: '%s' ", err.Error())
	}
	hash := md5.Sum(responseBytes)
	return &Result{MD5Hash: hex.EncodeToString(hash[:]), URLAddress: url.String()}, nil
}

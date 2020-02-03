package elasticsearch

import (
	"context"
	"net/http"

	rchttp "github.com/ONSdigital/dp-rchttp"
)

// ServiceName elasticsearch
const ServiceName = "elasticsearch"

//go:generate moq -out ./mock/rc-http.go -pkg mock . RchttpClient

// RchttpClient - interface representing a dp-rchttp client
type RchttpClient interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

// Client is an ElasticSearch client containing an HTTP client to contact the elasticsearch API.
type Client struct {
	httpCli      RchttpClient
	url          string
	serviceName  string
	signRequests bool
}

// NewClient returns a new initialised elasticsearch client with the default rchttp client
func NewClient(url string, signRequests bool, maxRetries int) *Client {
	httpClient := rchttp.NewClient()
	httpClient.SetMaxRetries(maxRetries)
	return NewClientWithHTTPClient(url, signRequests, httpClient)
}

// NewClientWithHTTPClient returns a new initialised elasticsearch client with the provided HTTP client
func NewClientWithHTTPClient(url string, signRequests bool, httpCli RchttpClient) *Client {
	return &Client{
		httpCli:      httpCli,
		url:          url,
		serviceName:  ServiceName,
		signRequests: signRequests,
	}
}

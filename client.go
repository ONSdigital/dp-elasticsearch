package elasticsearch

import (
	"context"
	"net/http"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
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
	Check        *health.Check
}

// NewClient returns a new initialized elasticsearch client with the default rchttp client
func NewClient(url string, signRequests bool, maxRetries int) *Client {
	httpClient := rchttp.NewClient()
	httpClient.SetMaxRetries(maxRetries)
	return NewClientWithHTTPClient(url, signRequests, httpClient)
}

// NewClientWithHTTPClient returns a new initialized elasticsearch client with the provided HTTP cilent
func NewClientWithHTTPClient(url string, signRequests bool, httpCli RchttpClient) *Client {

	// Initial Check struct
	check := &health.Check{Name: ServiceName}

	// Create Client
	return &Client{
		httpCli:      httpCli,
		url:          url,
		serviceName:  ServiceName,
		signRequests: signRequests,
		Check:        check,
	}
}

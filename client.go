package elasticsearch

import (
	dphttp "github.com/ONSdigital/dp-net/http"
)

// ServiceName elasticsearch
const ServiceName = "elasticsearch"

// Client is an ElasticSearch client containing an HTTP client to contact the elasticsearch API.
type Client struct {
	httpCli      dphttp.Clienter
	url          string
	serviceName  string
	signRequests bool
	indexes      []string
}

// NewClient returns a new initialised elasticsearch client with the default dp-net/http client
func NewClient(url string, signRequests bool, maxRetries int, indexes []string) *Client {
	httpClient := dphttp.NewClient()
	httpClient.SetMaxRetries(maxRetries)
	return NewClientWithHTTPClient(url, signRequests, httpClient, indexes)
}

// NewClientWithHTTPClient returns a new initialised elasticsearch client with the provided HTTP client
func NewClientWithHTTPClient(url string, signRequests bool, httpCli dphttp.Clienter, indexes []string) *Client {
	return &Client{
		httpCli:      httpCli,
		url:          url,
		serviceName:  ServiceName,
		signRequests: signRequests,
		indexes:      indexes,
	}
}

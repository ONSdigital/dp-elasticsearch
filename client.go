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
	indices      []string
}

// NewClient returns a new initialised elasticsearch client with the default dp-net/http client
func NewClient(url string, signRequests bool, maxRetries int, indices []string) *Client {
	httpClient := dphttp.NewClient()
	httpClient.SetMaxRetries(maxRetries)
	return NewClientWithHTTPClient(url, signRequests, httpClient, indices)
}

// NewClientWithHTTPClient returns a new initialised elasticsearch client with the provided HTTP client
func NewClientWithHTTPClient(url string, signRequests bool, httpCli dphttp.Clienter, indices []string) *Client {
	return &Client{
		httpCli:      httpCli,
		url:          url,
		serviceName:  ServiceName,
		signRequests: signRequests,
		indices:      indices,
	}
}

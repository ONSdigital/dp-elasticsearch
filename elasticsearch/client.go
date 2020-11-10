package elasticsearch

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	dphttp "github.com/ONSdigital/dp-net/http"
	"github.com/ONSdigital/log.go/log"
	awsauth "github.com/smartystreets/go-aws-auth"
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
func NewClient(url string, signRequests bool, maxRetries int, indexes ...string) *Client {
	httpClient := dphttp.NewClient()
	httpClient.SetMaxRetries(maxRetries)
	return NewClientWithHTTPClient(url, signRequests, httpClient, indexes...)
}

// NewClientWithHTTPClient returns a new initialised elasticsearch client with the provided HTTP client
func NewClientWithHTTPClient(url string, signRequests bool, httpCli dphttp.Clienter, indexes ...string) *Client {
	return &Client{
		httpCli:      httpCli,
		url:          url,
		serviceName:  ServiceName,
		signRequests: signRequests,
		indexes:      indexes,
	}
}

//CreateIndex creates an index in elasticsearch
func (cli *Client) CreateIndex(ctx context.Context, indexName string, indexSettings []byte) (int, error) {

	indexPath := cli.url + "/" + indexName
	_, status, err := cli.callElastic(ctx, indexPath, "PUT", indexSettings)
	if err != nil {
		return status, err
	}
	return status, nil
}

//DeleteIndex deletes an index in elasticsearch
func (cli *Client) DeleteIndex(ctx context.Context, indexName string) (int, error) {

	indexPath := cli.url + "/" + indexName
	_, status, err := cli.callElastic(ctx, indexPath, "DELETE", nil)
	if err != nil {
		return status, err
	}
	return status, nil
}

//AddDocument adds a JSON document to elasticsearch
func (cli *Client) AddDocument(ctx context.Context, indexName, documentType, documentID string, document []byte) (int, error) {

	documentPath := cli.url + "/" + indexName + "/" + documentType + "/" + documentID
	_, status, err := cli.callElastic(ctx, documentPath, "PUT", document)
	if err != nil {
		return status, err
	}
	return status, nil

}

// CallElastic builds a request to elasticsearch based on the method, path and payload
func (cli *Client) callElastic(ctx context.Context, path, method string, payload []byte) ([]byte, int, error) {
	logData := log.Data{"url": path, "method": method}

	URL, err := url.Parse(path)
	if err != nil {
		log.Event(ctx, "failed to create url for elastic call", log.ERROR, log.Error(err), logData)
		return nil, 0, err
	}
	path = URL.String()
	logData["url"] = path

	var req *http.Request

	if payload != nil {
		req, err = http.NewRequest(method, path, bytes.NewReader(payload))
		req.Header.Add("Content-type", "application/json")
		logData["payload"] = string(payload)
	} else {
		req, err = http.NewRequest(method, path, nil)
	}
	// check req, above, didn't error
	if err != nil {
		log.Event(ctx, "failed to create request for call to elastic", log.ERROR, log.Error(err), logData)
		return nil, 0, err
	}

	if cli.signRequests {
		awsauth.Sign(req)
	}

	resp, err := cli.httpCli.Do(ctx, req)
	if err != nil {
		log.Event(ctx, "failed to call elastic", log.ERROR, log.Error(err), logData)
		return nil, 0, err
	}
	defer resp.Body.Close()

	logData["http_code"] = resp.StatusCode

	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Event(ctx, "failed to read response body from call to elastic", log.ERROR, log.Error(err), logData)
		return nil, resp.StatusCode, err
	}
	logData["json_body"] = string(jsonBody)
	logData["status_code"] = resp.StatusCode

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		log.Event(ctx, "failed", log.ERROR, log.Error(ErrorUnexpectedStatusCode), logData)
		return nil, resp.StatusCode, ErrorUnexpectedStatusCode
	}

	return jsonBody, resp.StatusCode, nil
}

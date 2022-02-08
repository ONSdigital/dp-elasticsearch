package v2

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ONSdigital/dp-elasticsearch/v3/client"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	dphttp "github.com/ONSdigital/dp-net/http"
	"github.com/ONSdigital/log.go/v2/log"
)

// ServiceName elasticsearch
const (
	ServiceName = "elasticsearch"
)

// Client is an ElasticSearch client containing an HTTP client to contact the elasticsearch API.
type Client struct {
	httpCli     dphttp.Clienter
	url         string
	serviceName string
	indexes     []string
}

// NewClient returns a new initialised elasticsearch client with the default dp-net/http client
func NewClient(url string, maxRetries int, indexes ...string) *Client {
	httpClient := dphttp.NewClient()
	httpClient.SetMaxRetries(maxRetries)
	return NewClientWithHTTPClient(url, httpClient, indexes...)
}

func NewClientWithHTTPClient(url string, httpCli dphttp.Clienter, indexes ...string) *Client {
	cli := &Client{
		httpCli:     httpCli,
		url:         url,
		serviceName: ServiceName,
		indexes:     indexes,
	}

	// healthcheck client should not retry when calling a healthcheck endpoint,
	// append to current paths as to not change the client setup by service
	paths := cli.httpCli.GetPathsWithNoRetries()
	paths = append(paths, string(pathHealth))
	cli.httpCli.SetPathsWithNoRetries(paths)

	return cli
}

// GetIndices gets an index from elasticsearch
func (cli *Client) GetIndices(ctx context.Context, indexPatterns []string) (int, []byte, error) {

	indexPath := cli.url + "/" + strings.Join(indexPatterns, ",")
	body, status, err := cli.callElastic(ctx, indexPath, "GET", nil)
	if err != nil {
		return status, body, err
	}
	return status, body, nil
}

// CreateIndex creates an index in elasticsearch
func (cli *Client) CreateIndex(ctx context.Context, indexName string, indexSettings []byte) error {

	indexPath := cli.url + "/" + indexName
	_, status, err := cli.callElastic(ctx, indexPath, "PUT", indexSettings)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return errors.New(fmt.Sprintf("failed to create index '%s'", indexName))
	}
	return nil
}

func (cli *Client) DeleteIndices(ctx context.Context, indices []string) (int, error) {
	return cli.DeleteIndex(ctx, indices[0])
}

// DeleteIndex deletes an index in elasticsearch
func (cli *Client) DeleteIndex(ctx context.Context, indexName string) (int, error) {

	indexPath := cli.url + "/" + indexName
	_, status, err := cli.callElastic(ctx, indexPath, "DELETE", nil)
	if err != nil {
		return status, err
	}
	return status, nil
}

// AddDocument adds a JSON document to elasticsearch
func (cli *Client) AddDocument(ctx context.Context, indexName, documentID string, document []byte, options *client.AddDocumentOptions) error {
	documentType := "_doc"
	if options != nil && options.DocumentType != "" {
		documentType = options.DocumentType
	}
	documentPath := cli.url + "/" + indexName + "/" + documentType + "/" + documentID
	_, status, err := cli.callElastic(ctx, documentPath, "PUT", document)
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		return errors.New("unable to add document to elasticsearch")
	}
	return nil

}

func (cli *Client) UpdateAliases(ctx context.Context, alias, addIndices, removeIndices []string) error {
	return errors.New("function 'UpdateAliases' unsupported in this client")
}

// BulkUpdate uses an HTTP post request to submit data to Elastic Search
func (cli *Client) BulkUpdate(ctx context.Context, esDestIndex string, esDestURL string, bulk []byte) ([]byte, int, error) {

	uri := fmt.Sprintf("%s/%s/_bulk", esDestURL, esDestIndex)
	jsonBody, status, err := cli.callElastic(ctx, uri, "POST", bulk)
	if err != nil {
		log.Error(ctx, "error posting bulk request %s", err)
		return jsonBody, status, err
	}

	return jsonBody, status, err
}

// CallElastic builds a request to elasticsearch based on the method, path and payload
func (cli *Client) callElastic(ctx context.Context, path, method string, payload []byte) ([]byte, int, error) {

	logData := log.Data{
		"url":    path,
		"method": method,
	}
	URL, err := url.Parse(path)
	if err != nil {
		log.Error(ctx, "failed to create url for elastic call", err, logData)
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
		log.Error(ctx, "failed to create request for call to elastic", err, logData)
		return nil, 0, err
	}

	resp, err := cli.httpCli.Do(ctx, req)
	if err != nil {
		log.Error(ctx, "failed to call elastic", err, logData)
		return nil, 0, err
	}
	defer resp.Body.Close()

	logData["http_code"] = resp.StatusCode

	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(ctx, "failed to read response body from call to elastic", err, logData)
		return nil, resp.StatusCode, err
	}

	logData["json_body"] = string(jsonBody)
	logData["status_code"] = resp.StatusCode

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		log.Error(ctx, "failed as unexpected code", ErrorUnexpectedStatusCode, logData)
		return nil, resp.StatusCode, ErrorUnexpectedStatusCode
	}

	log.Info(ctx, "es response with response status code", logData)

	return jsonBody, resp.StatusCode, nil
}

func (cli *Client) BulkIndexAdd(ctx context.Context, indexName, documentID string, document []byte) error {
	return errors.New("bulk index add is not supported for legacy client")
}

// Close waits until all added items are flushed and closes the indexer.
func (cli *Client) BulkIndexClose(context.Context) error {
	return errors.New("bulk index close is currently not supported for legacy client")
}

package v2

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ONSdigital/dp-elasticsearch/v3/client"
	esError "github.com/ONSdigital/dp-elasticsearch/v3/errors"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
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
func NewClient(esURL string, maxRetries int, indexes ...string) *Client {
	httpClient := dphttp.NewClient()
	httpClient.SetMaxRetries(maxRetries)
	return NewClientWithHTTPClient(esURL, httpClient, indexes...)
}

func NewClientWithHTTPClient(esURL string, httpCli dphttp.Clienter, indexes ...string) *Client {
	cli := &Client{
		httpCli:     httpCli,
		url:         esURL,
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

// GetAlias returns a list of indices.
func (cli *Client) GetAlias(ctx context.Context) ([]byte, error) {
	return nil, errors.New("get alias is currently not supported by legacy client")
}

// GetIndices gets an index from elasticsearch
func (cli *Client) GetIndices(ctx context.Context, indexPatterns []string) (body []byte, err error) {
	indexPath := cli.url + "/" + strings.Join(indexPatterns, ",")

	body, status, err := cli.callElastic(ctx, indexPath, "GET", nil)
	if err != nil {
		return body, esError.StatusError{
			Err:  err,
			Code: status,
		}
	}

	return body, nil
}

// CreateIndex creates an index in elasticsearch
func (cli *Client) CreateIndex(ctx context.Context, indexName string, indexSettings []byte) error {
	indexPath := cli.url + "/" + indexName

	_, status, err := cli.callElastic(ctx, indexPath, "PUT", indexSettings)
	if err != nil {
		return esError.StatusError{
			Err:  err,
			Code: status,
		}
	}

	if status != http.StatusOK {
		return esError.StatusError{
			Err:  fmt.Errorf("failed to create index '%s'", indexName),
			Code: status,
		}
	}

	return nil
}

func (cli *Client) DeleteIndices(ctx context.Context, indices []string) error {
	return cli.DeleteIndex(ctx, indices[0])
}

func (cli *Client) Count(ctx context.Context, count client.Count) ([]byte, error) {
	return nil, errors.New("doc count is not supported for the legacy client")
}

// CountIndices feature is not supported for ES version 2.
func (cli *Client) CountIndices(ctx context.Context, indices []string) ([]byte, error) {
	return nil, errors.New("count is not supported for the legacy client")
}

// DeleteIndex deletes an index in elasticsearch
func (cli *Client) DeleteIndex(ctx context.Context, indexName string) error {
	indexPath := cli.url + "/" + indexName

	_, status, err := cli.callElastic(ctx, indexPath, "DELETE", nil)
	if err != nil {
		return esError.StatusError{
			Err:  err,
			Code: status,
		}
	}

	return nil
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
		return esError.StatusError{
			Err:  err,
			Code: status,
		}
	}

	if status != http.StatusCreated {
		return esError.StatusError{
			Err:  errors.New("unable to add document to elasticsearch"),
			Code: status,
		}
	}

	return nil
}

func (cli *Client) UpdateAliases(ctx context.Context, alias string, addIndices, removeIndices []string) error {
	return errors.New("function 'UpdateAliases' unsupported in this client")
}

// BulkUpdate uses an HTTP post request to submit data to Elastic Search
func (cli *Client) BulkUpdate(ctx context.Context, esDestIndex, esDestURL string, bulk []byte) ([]byte, error) {
	uri := fmt.Sprintf("%s/%s/_bulk", esDestURL, esDestIndex)

	body, status, err := cli.callElastic(ctx, uri, "POST", bulk)
	if err != nil {
		return body, esError.StatusError{
			Err:  err,
			Code: status,
		}
	}

	return body, nil
}

// CallElastic builds a request to elasticsearch based on the method, path and payload
func (cli *Client) callElastic(ctx context.Context, path, method string, payload []byte) (body []byte, status int, err error) {
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
		req, err = http.NewRequest(method, path, http.NoBody)
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

	jsonBody, err := io.ReadAll(resp.Body)
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

func (cli *Client) NewBulkIndexer(ctx context.Context) error {
	return errors.New("bulk indexer is not supported for legacy client")
}

func (cli *Client) BulkIndexAdd(ctx context.Context, action client.BulkIndexerAction, index, documentID string, document []byte, onSuccess client.SuccessFunc, onFailure client.FailureFunc) error {
	return errors.New("bulk index add is not supported for legacy client")
}

// Close waits until all added items are flushed and closes the indexer.
func (cli *Client) BulkIndexClose(context.Context) error {
	return errors.New("bulk index close is currently not supported for legacy client")
}

func (cli *Client) MultiSearch(ctx context.Context, searches []client.Search, queryParams *client.QueryParams) ([]byte, error) {
	return nil, errors.New("multi search is not supported for legacy client")
}

func (cli *Client) Search(ctx context.Context, search client.Search) ([]byte, error) {
	return nil, errors.New("search is not supported for legacy client")
}

func (cli *Client) Explain(ctx context.Context, documentID string, search client.Search) ([]byte, error) {
	return nil, errors.New("explain is not supported for legacy client")
}

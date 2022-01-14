package go_elasticsearch_v710

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	es710 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

const (
	bulkIndexerClientShouldNotBeNilErrMsg = "bulk indexer client should not be nil"
)

type ESClient struct {
	esClient    *es710.Client
	bulkIndexer *bulkIndexer
}

// NewESClient returns a new elastic search client version 7.10
func NewESClient(rawURL string, transport http.RoundTripper, indexName string) (*ESClient, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, errors.New("failed to specify valid elasticsearch url")
	}
	if indexName == "" {
		return nil, errors.New("should specify a valid index name")
	}
	newESClient, err := es710.NewClient(es710.Config{
		Addresses: []string{parsedURL.String()},
		Transport: transport,
	})
	if err != nil {
		return nil, err
	}
	bi, err := newBulkIndexer(indexName, newESClient)
	if err != nil {
		return nil, err
	}
	return &ESClient{
		esClient:    newESClient,
		bulkIndexer: bi,
	}, nil
}

// GetIndices  returns information about one or more indices.
func (cli *ESClient) GetIndices(ctx context.Context, indexPatterns []string) (int, []byte, error) {
	res, err := cli.esClient.Indices.Get(indexPatterns)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return res.StatusCode, nil, errors.New("error occured while trying to retrieve indices")
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, data, nil
}

// IndicesCreate creates an index with optional settings and mappings.
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/master/indices-create-index.html.
func (cli *ESClient) CreateIndex(ctx context.Context, indexName string, indexSettings []byte) (int, error) {
	res, err := cli.esClient.Indices.Create(indexName, cli.esClient.Indices.Create.WithBody(bytes.NewReader(indexSettings)))
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return res.StatusCode, errors.New("error occured while trying to create index")
	}
	return res.StatusCode, nil
}

// IndicesDelete deletes an index.
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/master/indices-delete-index.html.
func (cli *ESClient) DeleteIndex(ctx context.Context, indexName string) (int, error) {
	return cli.DeleteIndices(ctx, []string{indexName})
}

// IndicesDelete deletes an index.
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/master/indices-delete-index.html.
func (cli *ESClient) DeleteIndices(ctx context.Context, indices []string) (int, error) {
	res, err := cli.esClient.Indices.Delete(indices)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return res.StatusCode, errors.New("error occured while trying to create index")
	}
	return res.StatusCode, nil
}

func (cli *ESClient) AddDocument(ctx context.Context, indexName, documentType, documentID string, document []byte) (int, error) {
	res, err := esapi.CreateRequest{
		Index:        indexName,
		DocumentID:   documentID,
		Body:         bytes.NewReader(document),
		DocumentType: documentType,
	}.Do(ctx, cli.esClient)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return res.StatusCode, errors.New("error occured while trying to add document")
	}
	return res.StatusCode, nil
}

// Bulk allows to perform multiple index/update/delete operations in a single request.
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-bulk.html.
func (cli *ESClient) BulkUpdate(ctx context.Context, indexName, url string, payload []byte) ([]byte, int, error) {
	res, err := esapi.BulkRequest{
		Index: indexName,
		Body:  bytes.NewReader(payload),
	}.Do(ctx, cli.esClient)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, res.StatusCode, errors.New("error occured while trying to bulk update document")
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, err
	}
	return data, res.StatusCode, nil
}

// Add adds an item to the indexer. It returns an error when the item cannot be added.
// Use the OnSuccess and OnFailure callbacks to get the operation result for the item.
//
// You must call the Close() method after you're done adding items.
//
// It is safe for concurrent use. When it's called from goroutines,
// they must finish before the call to Close, eg. using sync.WaitGroup.
func (cli *ESClient) BulkIndexAdd(ctx context.Context, indexName, documentID string, document []byte) error {
	return cli.bulkIndexer.Add(ctx, documentID, document)
}

// Close waits until all added items are flushed and closes the indexer.
func (cli *ESClient) BulkIndexClose(ctx context.Context) error {
	return cli.bulkIndexer.Close(ctx)
}

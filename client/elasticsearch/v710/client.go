package v710

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/ONSdigital/dp-elasticsearch/v3/client"

	es710 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

const (
	bulkIndexerClientShouldNotBeNilErrMsg = "bulk indexer client should not be nil"
)

type ESClient struct {
	bulkIndexer *bulkIndexer
	esClient    *es710.Client
	indexes     []string
}

// NewESClient returns a new elastic search client version 7.10
func NewESClient(esURL string, transport http.RoundTripper) (*ESClient, error) {
	parsedURL, err := url.ParseRequestURI(esURL)
	if err != nil {
		return nil, errors.New("failed to specify valid elasticsearch url")
	}

	newESClient, err := es710.NewClient(es710.Config{
		Addresses: []string{parsedURL.String()},
		Transport: transport,
	})
	if err != nil {
		return nil, err
	}

	return &ESClient{
		esClient: newESClient,
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
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/7.10/indices-create-index.html.
func (cli *ESClient) CreateIndex(ctx context.Context, indexName string, indexSettings []byte) error {
	res, err := cli.esClient.Indices.Create(indexName, cli.esClient.Indices.Create.WithBody(bytes.NewReader(indexSettings)))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New("error occured while trying to create index")
	}

	return nil
}

// IndicesDelete deletes an index.
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/7.10/indices-delete-index.html.
func (cli *ESClient) DeleteIndex(ctx context.Context, indexName string) (int, error) {
	return cli.DeleteIndices(ctx, []string{indexName})
}

// IndicesDelete deletes an index.
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/7.10/indices-delete-index.html.
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

// AddDocument adds a document to the index specified. Upsert option not implemented.
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/7.10/docs-update.html.
func (cli *ESClient) AddDocument(ctx context.Context, indexName, documentID string, document []byte, options *client.AddDocumentOptions) error {
	req := esapi.CreateRequest{
		Index:      indexName,
		DocumentID: documentID,
		Body:       bytes.NewReader(document),
	}

	if options != nil {
		if options.DocumentType != "" {
			req.DocumentType = options.DocumentType
		}

		if options.Upsert {
			return errors.New("es710 client currently cannot handle upsert option when creating a document")
		}
	}

	res, err := req.Do(ctx, cli.esClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New("error occured while trying to add document")
	}

	return nil
}

// UpdateAliases removes and adds an alias to indexes.
func (cli *ESClient) UpdateAliases(ctx context.Context, alias string, removeIndices, addIndices []string) error {
	var actions []string

	if len(removeIndices) > 0 {
		removeAction := fmt.Sprintf(
			`{"remove": {"indices": "%s","alias": "%s"}}`,
			strings.Join(removeIndices, `","`),
			alias)
		actions = append(actions, removeAction)
	}

	if len(addIndices) > 0 {
		addAction := fmt.Sprintf(
			`{"add": {"indices": "%s","alias": "%s"}}`,
			strings.Join(addIndices, `","`),
			alias)
		actions = append(actions, addAction)
	}

	update := fmt.Sprintf(
		`{"actions": [%s]}`,
		strings.Join(actions, ","))
	_, err := cli.esClient.Indices.UpdateAliases(strings.NewReader(update))
	if err != nil {
		return err
	}

	return nil
}

// Bulk allows to perform multiple index/update/delete operations in a single request.
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/7.10/docs-bulk.html.
func (cli *ESClient) BulkUpdate(ctx context.Context, indexName, esUrl string, payload []byte) ([]byte, int, error) {
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
func (cli *ESClient) BulkIndexAdd(ctx context.Context, action client.BulkIndexerAction, index, documentID string, document []byte) error {
	return cli.bulkIndexer.Add(ctx, action, index, documentID, document)
}

// Close waits until all added items are flushed and closes the indexer.
func (cli *ESClient) BulkIndexClose(ctx context.Context) error {
	return cli.bulkIndexer.Close(ctx)
}
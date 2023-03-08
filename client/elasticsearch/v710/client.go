package v710

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ONSdigital/dp-elasticsearch/v3/client"
	esError "github.com/ONSdigital/dp-elasticsearch/v3/errors"
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

// GetAlias returns a list of indices.
func (cli *ESClient) GetAlias(ctx context.Context) ([]byte, error) {
	res, err := cli.esClient.Indices.GetAlias()
	if err != nil {
		return nil, esError.StatusError{Err: err, Code: getStatusCode(res)}
	}
	defer res.Body.Close()

	if err = checkForError(res); err != nil {
		return nil, esError.StatusError{
			Err:  fmt.Errorf("error occured while trying to retrieve aliases: %w", err),
			Code: getStatusCode(res),
		}
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, esError.StatusError{
			Err:  err,
			Code: getStatusCode(res),
		}
	}

	return data, nil
}

// GetIndices  returns information about one or more indices.
func (cli *ESClient) GetIndices(ctx context.Context, indexPatterns []string) ([]byte, error) {
	res, err := cli.esClient.Indices.Get(indexPatterns)
	if err != nil {
		return nil, esError.StatusError{Err: err, Code: getStatusCode(res)}
	}
	defer res.Body.Close()

	if err = checkForError(res); err != nil {
		return nil, esError.StatusError{
			Err:  fmt.Errorf("error occured while trying to retrieve indices: %w", err),
			Code: getStatusCode(res),
		}
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, esError.StatusError{
			Err:  err,
			Code: getStatusCode(res),
		}
	}

	return data, nil
}

// IndicesCreate creates an index with optional settings and mappings.
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/7.10/indices-create-index.html.
func (cli *ESClient) CreateIndex(ctx context.Context, indexName string, indexSettings []byte) error {
	res, err := cli.esClient.Indices.Create(indexName, cli.esClient.Indices.Create.WithBody(bytes.NewReader(indexSettings)))
	if err != nil {
		return esError.StatusError{
			Err:  err,
			Code: getStatusCode(res),
		}
	}
	defer res.Body.Close()

	if err := checkForError(res); err != nil {
		return esError.StatusError{
			Err:  fmt.Errorf("error occured while trying to create index: %w", err),
			Code: getStatusCode(res),
		}
	}

	return nil
}

// IndicesDelete deletes an index.
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/7.10/indices-delete-index.html.
func (cli *ESClient) DeleteIndex(ctx context.Context, indexName string) error {
	return cli.DeleteIndices(ctx, []string{indexName})
}

// IndicesDelete deletes an index.
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/7.10/indices-delete-index.html.
func (cli *ESClient) DeleteIndices(ctx context.Context, indices []string) error {
	res, err := cli.esClient.Indices.Delete(indices)
	if err != nil {
		return esError.StatusError{
			Err:  err,
			Code: getStatusCode(res),
		}
	}
	defer res.Body.Close()

	if err := checkForError(res); err != nil {
		return esError.StatusError{
			Err:  fmt.Errorf("error occured while trying to delete index: %w", err),
			Code: getStatusCode(res),
		}
	}

	return nil
}

// Count returns number of documents matching a query.
//
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/master/search-count.html.
func (cli *ESClient) Count(ctx context.Context, count client.Count) ([]byte, error) {
	countReq := func(r *esapi.CountRequest) {
		r.Body = bytes.NewReader(count.Query)
	}
	res, err := cli.esClient.Count(countReq)
	if err != nil {
		return nil, esError.StatusError{
			Err:  err,
			Code: getStatusCode(res),
		}
	}
	defer res.Body.Close()

	if err = checkForError(res); err != nil {
		return nil, esError.StatusError{
			Err:  fmt.Errorf("error occured while trying to count indicies: %w", err),
			Code: getStatusCode(res),
		}
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, esError.StatusError{
			Err:  err,
			Code: getStatusCode(res),
		}
	}

	return data, nil
}

// Count returns number of documents matching a query.
//
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/master/search-count.html.
func (cli *ESClient) CountIndices(ctx context.Context, indices []string) ([]byte, error) {
	countReq := func(r *esapi.CountRequest) {
		r.Index = indices
	}
	res, err := cli.esClient.Count(countReq)
	if err != nil {
		return nil, esError.StatusError{
			Err:  err,
			Code: getStatusCode(res),
		}
	}
	defer res.Body.Close()

	if err = checkForError(res); err != nil {
		return nil, esError.StatusError{
			Err:  fmt.Errorf("error occured while trying to count indices: %w", err),
			Code: getStatusCode(res),
		}
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, esError.StatusError{
			Err:  err,
			Code: getStatusCode(res),
		}
	}

	return data, nil
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
			return esError.StatusError{
				Err: errors.New("es710 client currently cannot handle upsert option when creating a document"),
			}
		}
	}

	res, err := req.Do(ctx, cli.esClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if err := checkForError(res); err != nil {
		return esError.StatusError{
			Err:  fmt.Errorf("error occured while trying to add document: %w", err),
			Code: getStatusCode(res),
		}
	}

	return nil
}

func (cli *ESClient) Explain(ctx context.Context, documentID string, search client.Search) ([]byte, error) {
	req := esapi.ExplainRequest{
		Index:      search.Header.Index,
		DocumentID: documentID,
		Body:       bytes.NewReader(search.Query),
	}

	res, err := req.Do(ctx, cli.esClient)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err = checkForError(res); err != nil {
		return nil, esError.StatusError{
			Err:  fmt.Errorf("error occured while trying to call explain api: %w", err),
			Code: getStatusCode(res),
		}
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, esError.StatusError{
			Err:  err,
			Code: getStatusCode(res),
		}
	}
	return data, nil
}

// Search returns results matching a query.
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/master/search-search.html.
func (cli *ESClient) Search(ctx context.Context, search client.Search) ([]byte, error) {
	req := esapi.SearchRequest{
		Index: []string{search.Header.Index},
		Body:  bytes.NewReader(search.Query),
	}

	res, err := req.Do(ctx, cli.esClient)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err = checkForError(res); err != nil {
		return nil, esError.StatusError{
			Err:  fmt.Errorf("error occured while trying to search documents: %w", err),
			Code: getStatusCode(res),
		}
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, esError.StatusError{
			Err:  err,
			Code: getStatusCode(res),
		}
	}

	return data, nil
}

// Msearch allows to execute several search operations in one request.
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/master/search-multi-search.html.
func (cli *ESClient) MultiSearch(ctx context.Context, searches []client.Search, queryParams *client.QueryParams) ([]byte, error) {
	body, err := convertToMultilineSearches(searches)
	if err != nil {
		return nil, err
	}
	req := esapi.MsearchRequest{
		Body: bytes.NewReader(body),
	}
	if queryParams != nil && queryParams.EnableTotalHitsCounter != nil {
		req.RestTotalHitsAsInt = queryParams.EnableTotalHitsCounter
	}

	res, err := req.Do(ctx, cli.esClient)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err = checkForError(res); err != nil {
		return nil, esError.StatusError{
			Err:  fmt.Errorf("error occured while trying to multi search documents: %w", err),
			Code: getStatusCode(res),
		}
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, esError.StatusError{
			Err:  err,
			Code: getStatusCode(res),
		}
	}

	return data, nil
}

// UpdateAliases removes and adds an alias to indexes.
func (cli *ESClient) UpdateAliases(ctx context.Context, alias string, removeIndices, addIndices []string) error {
	var actions []string

	if len(removeIndices) > 0 {
		removeAction := fmt.Sprintf(
			`{"remove": {"indices": %q,"alias": %q}}`,
			strings.Join(removeIndices, `","`),
			alias)
		actions = append(actions, removeAction)
	}

	if len(addIndices) > 0 {
		addAction := fmt.Sprintf(
			`{"add": {"indices": %q,"alias": %q}}`,
			strings.Join(addIndices, `","`),
			alias)
		actions = append(actions, addAction)
	}

	update := fmt.Sprintf(
		`{"actions": [%s]}`,
		strings.Join(actions, ","))
	res, err := cli.esClient.Indices.UpdateAliases(strings.NewReader(update))
	if err != nil {
		return esError.StatusError{
			Err:  err,
			Code: getStatusCode(res),
		}
	}

	return nil
}

// Bulk allows to perform multiple index/update/delete operations in a single request.
// See full documentation at https://www.elastic.co/guide/en/elasticsearch/reference/7.10/docs-bulk.html.
func (cli *ESClient) BulkUpdate(ctx context.Context, indexName, esURL string, payload []byte) ([]byte, error) {
	res, err := esapi.BulkRequest{
		Index: indexName,
		Body:  bytes.NewReader(payload),
	}.Do(ctx, cli.esClient)
	if err != nil {
		return nil, esError.StatusError{
			Err:  err,
			Code: getStatusCode(res),
		}
	}
	defer res.Body.Close()

	if err = checkForError(res); err != nil {
		return nil, esError.StatusError{
			Err:  fmt.Errorf("error occured while trying to bulk update document: %w", err),
			Code: getStatusCode(res),
		}
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, esError.StatusError{
			Err:  err,
			Code: getStatusCode(res),
		}
	}

	return data, nil
}

// NewBulkIndexer creates a bulkIndexer for use of the client.
func (cli *ESClient) NewBulkIndexer(ctx context.Context) error {
	bulkIndexer, err := newBulkIndexer(cli.esClient)
	if err != nil {
		return esError.StatusError{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}

	cli.bulkIndexer = bulkIndexer

	return nil
}

// Add adds an item to the indexer. It returns an error when the item cannot be added.
// Use the OnSuccess and OnFailure callbacks to get the operation result for the item.
//
// You must call the Close() method after you're done adding items.
//
// It is safe for concurrent use. When it's called from goroutines,
// they must finish before the call to Close, eg. using sync.WaitGroup.
func (cli *ESClient) BulkIndexAdd(
	ctx context.Context,
	action client.BulkIndexerAction,
	index,
	documentID string,
	document []byte,
	onSuccess client.SuccessFunc,
	onFailure client.FailureFunc,
) error {
	if cli.bulkIndexer == nil {
		return esError.StatusError{
			Err:  errors.New(bulkIndexerClientShouldNotBeNilErrMsg),
			Code: http.StatusInternalServerError,
		}
	}

	return cli.bulkIndexer.Add(ctx, action, index, documentID, document, onSuccess, onFailure)
}

// Close waits until all added items are flushed and closes the indexer.
func (cli *ESClient) BulkIndexClose(ctx context.Context) error {
	if cli.bulkIndexer == nil {
		return esError.StatusError{
			Err:  errors.New(bulkIndexerClientShouldNotBeNilErrMsg),
			Code: http.StatusInternalServerError,
		}
	}

	return cli.bulkIndexer.Close(ctx)
}

func convertToMultilineSearches(searches []client.Search) (body []byte, err error) {
	for _, search := range searches {
		headerByte, err := json.Marshal(search.Header)
		if err != nil {
			return nil, err
		}
		body = append(body, headerByte...)
		body = append(body, '\n')
		body = append(body, search.Query...)
		body = append(body, '\n')
	}
	return body, nil
}

// getStatusCode returns the response StatusCode, or 0 if res is nil
func getStatusCode(res *esapi.Response) int {
	if res == nil {
		return 0
	}
	return res.StatusCode
}

// checkForError checks if the provided elasticsearch response contains an error.
// if it does, it is read and returned as a string error
func checkForError(res *esapi.Response) error {
	if res == nil {
		return errors.New("nil elasticsearch api response")
	}

	if !res.IsError() {
		return nil
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to ready elasticsearch response body for an error case: %w", err)
	}

	return fmt.Errorf("error response from elasticsearch: %s", string(resBody))
}

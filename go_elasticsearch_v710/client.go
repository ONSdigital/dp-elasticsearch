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

type ESClient struct {
	esClient *es710.Client
}

func NewESClient(rawURL string, transport http.RoundTripper) (*ESClient, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
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

// GetIndices gets an index from elasticsearch
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

func (cli *ESClient) DeleteIndex(ctx context.Context, indexName string) (int, error) {
	return cli.DeleteIndices(ctx, []string{indexName})
}

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

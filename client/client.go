package client

import (
	"context"
	"net/http"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

//go:generate moq -out ./mocks/client.go -pkg mocks . Client

// Client holds the methods for ElasticSearch clients
type Client interface {
	AddDocument(ctx context.Context, indexName, documentID string, document []byte, opts *AddDocumentOptions) error
	BulkUpdate(ctx context.Context, indexName, url string, settings []byte) ([]byte, error)
	BulkIndexAdd(ctx context.Context, action BulkIndexerAction, index, documentID string, document []byte, onSuccess SuccessFunc, onFailure FailureFunc) error
	BulkIndexClose(context.Context) error
	Checker(ctx context.Context, state *health.CheckState) error
	CreateIndex(ctx context.Context, indexName string, indexSettings []byte) error
	DeleteDocument(ctx context.Context, indexName, documentID string) error
	DeleteDocumentByQuery(ctx context.Context, search Search) error
	DeleteIndex(ctx context.Context, indexName string) error
	DeleteIndices(ctx context.Context, indices []string) error
	GetAlias(ctx context.Context) ([]byte, error)
	GetIndices(ctx context.Context, indexPatterns []string) ([]byte, error)
	NewBulkIndexer(context.Context) error
	UpdateAliases(ctx context.Context, alias string, removeIndices, addIndices []string) error
	MultiSearch(ctx context.Context, searches []Search, queryParams *QueryParams) ([]byte, error)
	Search(ctx context.Context, search Search) ([]byte, error)
	CountIndices(ctx context.Context, indices []string) ([]byte, error)
	Count(ctx context.Context, count Count) ([]byte, error)
	Explain(ctx context.Context, documentID string, search Search) ([]byte, error)
}

type Library string

type BulkIndexerAction string

const (
	GoElasticV710 Library = "GoElastic_v710"
	OpenSearch    Library = "OpenSearch"
)

// Config holds the configuration of search client
type Config struct {
	ClientLib  Library
	MaxRetries int
	Address    string
	Indexes    []string
	Transport  http.RoundTripper
}

type AddDocumentOptions struct {
	DocumentType string // Deprecated - not used by newer versions of elasticsearch
	Upsert       bool
}

type Header struct {
	Index string `json:"index"`
}

type Search struct {
	Header Header
	Query  []byte
}

type Count struct {
	Query []byte
}

type QueryParams struct {
	EnableTotalHitsCounter *bool
}

// SuccessFunc is the callback func signature for a successful bulk add operation, as expected by go-elasticsearch
type SuccessFunc = func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem)

// FailureFunc is the callback func signature for a successful bulk add operation, as expected by go-elasticsearch
type FailureFunc = func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error)

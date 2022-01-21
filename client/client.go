package client

import (
	"context"
	"net/http"
)

type Client interface {
	GetIndices(ctx context.Context, indexPatterns []string) (int, []byte, error)
	CreateIndex(ctx context.Context, indexName string, indexSettings []byte) error
	DeleteIndex(ctx context.Context, indexName string) (int, error)
	DeleteIndices(ctx context.Context, indices []string) (int, error)
	AddDocument(ctx context.Context, indexName, documentID string, document []byte, opts *AddDocumentOptions) error
	BulkUpdate(ctx context.Context, indexName, url string, settings []byte) ([]byte, int, error)
	BulkIndexAdd(ctx context.Context, indexName, documentID string, document []byte) error
	BulkIndexClose(context.Context) error
}

type ClientLibrary string

const (
	GoElastic_V710 ClientLibrary = "GoElastic_v710"
	OpenSearch     ClientLibrary = "OpenSearch"
)

// Config holds the configuration of search client
type Config struct {
	ClientLib  ClientLibrary
	MaxRetries int
	Address    string
	Indexes    []string
	Transport  http.RoundTripper
}

type AddDocumentOptions struct {
	DocumentType string // Deprecated - not used by newer versions of elasticsearch
	Upsert bool
}


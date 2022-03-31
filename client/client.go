package client

import (
	"context"
	"net/http"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
)

//go:generate moq -out ./mocks/client.go -pkg mocks . Client

type Client interface {
	AddDocument(ctx context.Context, indexName, documentID string, document []byte, opts *AddDocumentOptions) error
	BulkUpdate(ctx context.Context, indexName, url string, settings []byte) ([]byte, error)
	BulkIndexAdd(ctx context.Context, action BulkIndexerAction, index, documentID string, document []byte) error
	BulkIndexClose(context.Context) error
	Checker(ctx context.Context, state *health.CheckState) error
	CreateIndex(ctx context.Context, indexName string, indexSettings []byte) error
	DeleteIndex(ctx context.Context, indexName string) error
	DeleteIndices(ctx context.Context, indices []string) error
	GetAlias(ctx context.Context) ([]byte, error)
	GetIndices(ctx context.Context, indexPatterns []string) ([]byte, error)
	NewBulkIndexer(context.Context) error
	UpdateAliases(ctx context.Context, alias string, removeIndices, addIndices []string) error
	MultiSearch(ctx context.Context, indices []string, document []byte) error
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

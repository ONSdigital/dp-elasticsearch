package elasticsearch

import (
	"context"
	"fmt"
	"github.com/ONSdigital/dp-elasticsearch/v3/clients/elasticsearch/v2"
	"github.com/ONSdigital/dp-elasticsearch/v3/clients/elasticsearch/v710"
	"net/http"
)

type ClientLibrary string

const (
	GoElastic_V710 ClientLibrary = "GoElastic_v710"
	OpenSearch     ClientLibrary = "OpenSearch"
)

type Client interface {
	GetIndices(ctx context.Context, indexPatterns []string) (int, []byte, error)
	CreateIndex(ctx context.Context, indexName string, indexSettings []byte)  error
	DeleteIndex(ctx context.Context, indexName string) (int, error)
	DeleteIndices(ctx context.Context, indices []string) (int, error)
	AddDocument(ctx context.Context, indexName, documentType, documentID string, document []byte) (int, error)
	BulkUpdate(ctx context.Context, indexName, url string, settings []byte) ([]byte, int, error)
	BulkIndexAdd(ctx context.Context, indexName, documentID string, document []byte) error
	BulkIndexClose(context.Context) error
}

// Config holds the configuration of search client
type Config struct {
	ClientLib  ClientLibrary
	MaxRetries int
	Address    string
	Indexes    []string
	Transport  http.RoundTripper
}

func NewClient(cfg Config) (Client, error) {
	switch cfg.ClientLib {
	case GoElastic_V710:
		return v710.NewESClient(cfg.Address, cfg.Transport)
	case OpenSearch:
		return nil, fmt.Errorf("Opensearch client is currently not implemented")
	default:
		return v2.NewClient(cfg.Address, cfg.MaxRetries, cfg.Indexes...), nil
	}
}

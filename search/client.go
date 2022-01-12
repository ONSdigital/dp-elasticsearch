package search

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-elasticsearch/v2/elasticsearch"
	"github.com/ONSdigital/dp-elasticsearch/v2/go_elasticsearch_v710"
)

type ClientLibrary string

const (
	GoElastic_V710 ClientLibrary = "GoElastic_v710"
	OpenSearch     ClientLibrary = "OpenSearch"
)

type Client interface {
	GetIndices(ctx context.Context, indexPatterns []string) (int, []byte, error)
	CreateIndex(ctx context.Context, indexName string, indexSettings []byte) (int, error)
	DeleteIndex(ctx context.Context, indexName string) (int, error)
	DeleteIndices(ctx context.Context, indices []string) (int, error)
	AddDocument(ctx context.Context, indexName, documentType, documentID string, document []byte) (int, error)
	BulkUpdate(ctx context.Context, indexName, url string, settings []byte) ([]byte, int, error)
}

// Config holds the configuration of search client
type Config struct {
	ClientLib  ClientLibrary
	MaxRetries int
	Address    string
	indexes    []string
	Transport  http.RoundTripper
}

func NewClient(cfg Config) (Client, error) {
	switch cfg.ClientLib {
	case GoElastic_V710:
		return go_elasticsearch_v710.NewESClient(cfg.Address, cfg.Transport)
	case OpenSearch:
		return nil, fmt.Errorf("Opensearch client is currently not implemented")
	default:
		return elasticsearch.NewClient(cfg.Address, cfg.MaxRetries, cfg.indexes...), nil
	}
	return nil, nil
}

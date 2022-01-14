package go_elasticsearch_v710

import (
	"bytes"
	"context"
	"errors"
	"time"

	es710 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

const numWorkers = 5

type bulkIndexer struct {
	bi esutil.BulkIndexer
}

// NewBulkIndexer creates a new bulk indexer.
func newBulkIndexer(indexName string, es *es710.Client) (*bulkIndexer, error) {
	if indexName == "" {
		return nil, errors.New("index name should not be empty")
	}
	if es == nil {
		return nil, errors.New("elastic client should not be nil")
	}
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         indexName,
		Client:        es,
		NumWorkers:    numWorkers,
		FlushInterval: 30 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &bulkIndexer{
		bi: bi,
	}, nil
}

// Add adds an item to the indexer. It returns an error when the item cannot be added.
// Use the OnSuccess and OnFailure callbacks to get the operation result for the item.
//
// You must call the Close() method after you're done adding items.
//
// It is safe for concurrent use. When it's called from goroutines,
// they must finish before the call to Close, eg. using sync.WaitGroup.
func (b *bulkIndexer) Add(ctx context.Context, documentID string, document []byte) error {
	bulkIndexerItem := esutil.BulkIndexerItem{
		Action:     "index",
		DocumentID: documentID,
		Body:       bytes.NewReader(document),
	}
	return b.bi.Add(ctx, bulkIndexerItem)
}

// Close waits until all added items are flushed and closes the indexer.
func (b *bulkIndexer) Close(ctx context.Context) error {
	return b.bi.Close(ctx)
}

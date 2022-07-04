package v710

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/ONSdigital/dp-elasticsearch/v4/client"
	es710 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

const numWorkers = 5

const (
	Create = client.BulkIndexerAction("create")
	Delete = client.BulkIndexerAction("delete")
	Index  = client.BulkIndexerAction("index")
	Update = client.BulkIndexerAction("update")
)

type bulkIndexer struct {
	bi esutil.BulkIndexer
}

// NewBulkIndexer creates a new bulk indexer.
func newBulkIndexer(es *es710.Client) (*bulkIndexer, error) {
	if es == nil {
		return nil, errors.New("elastic client should not be nil")
	}

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:        es,
		FlushInterval: 30 * time.Second,
		NumWorkers:    numWorkers,
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
func (b *bulkIndexer) Add(ctx context.Context, action client.BulkIndexerAction, index, documentID string, document []byte) error {
	bulkIndexerItem := esutil.BulkIndexerItem{
		Action:     string(action),
		Body:       bytes.NewReader(document),
		DocumentID: documentID,
		Index:      index,
	}

	return b.bi.Add(ctx, bulkIndexerItem)
}

// Close waits until all added items are flushed and closes the indexer.
func (b *bulkIndexer) Close(ctx context.Context) error {
	return b.bi.Close(ctx)
}

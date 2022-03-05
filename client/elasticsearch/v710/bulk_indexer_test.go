package v710

import (
	"context"
	"errors"
	"testing"

	es710 "github.com/elastic/go-elasticsearch/v7"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewBulkIndexer(t *testing.T) {

	Convey("Given non nil es710 client", t, func() {
		client := &es710.Client{}

		Convey("When calling newBulkIndexer", func() {
			expectedBulkIndexer := &bulkIndexer{}

			bulkIndexer, err := newBulkIndexer(client)

			Convey("Then a new bulk indexer is returned", func() {
				So(err, ShouldBeNil)
				So(bulkIndexer, ShouldHaveSameTypeAs, expectedBulkIndexer)
			})
		})
	})

	Convey("Given client is nil", t, func() {
		var client *es710.Client

		Convey("When calling newBulkIndexer", func() {
			bulkIndexer, err := newBulkIndexer(client)

			Convey("Then an error is returned", func() {
				So(err, ShouldResemble, errors.New("elastic client should not be nil"))
				So(bulkIndexer, ShouldBeNil)
			})
		})
	})
}

func TestBulkIndexerMethods(t *testing.T) {
	testCtx := context.Background()
	indexName := "test123"

	Convey("Given a valid bulk indexer", t, func() {
		bulkIndexer, err := setupBulkIndexer()
		if err != nil {
			t.Errorf("failed to setup bulk indexer for test")
		}

		Convey("When calling Add method", func() {
			err := bulkIndexer.Add(testCtx, Create, indexName, "123", []byte{})

			Convey("Then item is added to the bulk indexer without errors", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When calling Close method", func() {
			err := bulkIndexer.Close(testCtx)

			Convey("Then the bulkindexer was closed without errors", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func setupBulkIndexer() (*bulkIndexer, error) {
	return newBulkIndexer(&es710.Client{})
}

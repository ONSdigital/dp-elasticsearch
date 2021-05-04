package elasticsearch_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-elasticsearch/v2/elasticsearch"
	dphttp "github.com/ONSdigital/dp-net/http"
	. "github.com/smartystreets/goconvey/convey"
)

const testType = "_type"
const testID = "id"

var (
	errorUnexpectedStatusCode = errors.New("unexpected status code from api")
)

var doSuccessful = func(ctx context.Context, request *http.Request) (*http.Response, error) {
	return resp("do successful", 200), nil
}

var doUnsuccessful = func(ctx context.Context, request *http.Request) (*http.Response, error) {
	return resp("do unsuccessful", 0), ErrUnreachable
}

var unexpectedStatusCode = func(ctx context.Context, request *http.Request) (*http.Response, error) {
	return resp("unexpected status", 400), nil
}

var emptyListOfPathsWithNoRetries = func() []string {
	return []string{}
}

var setListOfPathsWithNoRetries = func(listOfPaths []string) {
	return
}

func clientMock(doFunc func(ctx context.Context, request *http.Request) (*http.Response, error)) *dphttp.ClienterMock {
	return &dphttp.ClienterMock{
		DoFunc:                    doFunc,
		GetPathsWithNoRetriesFunc: emptyListOfPathsWithNoRetries,
		SetPathsWithNoRetriesFunc: setListOfPathsWithNoRetries,
	}
}

func TestCreateIndex(t *testing.T) {

	indexSettings := []byte("settings")

	Convey("Given that an index with settings is created", t, func() {

		httpCli := clientMock(doSuccessful)
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)
		checkClient(httpCli)

		Convey("A status code of 200 and no error is returned", func() {
			status, err := cli.CreateIndex(context.Background(), testIndex, indexSettings)
			So(err, ShouldEqual, nil)
			So(httpCli.DoCalls(), ShouldHaveLength, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one")
			So(status, ShouldEqual, 200)
		})
	})

	Convey("Given that an index without settings is created", t, func() {

		httpCli := clientMock(doSuccessful)
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)
		checkClient(httpCli)

		Convey("A status code of 200 and no error is returned", func() {
			status, err := cli.CreateIndex(context.Background(), testIndex, nil)
			So(err, ShouldEqual, nil)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one")
			So(status, ShouldEqual, 200)
		})
	})

	Convey("Given that there is a server error", t, func() {

		httpCli := clientMock(doUnsuccessful)
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)
		checkClient(httpCli)

		Convey("A status code of 500 and an error is returned", func() {
			status, err := cli.CreateIndex(context.Background(), testIndex, indexSettings)
			So(err, ShouldNotEqual, nil)
			So(err, ShouldResemble, ErrUnreachable)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one")
			So(status, ShouldEqual, 0)
		})
	})

	Convey("Given that an elasticsearch returns an unexpected status code", t, func() {

		httpCli := clientMock(unexpectedStatusCode)
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)
		checkClient(httpCli)

		Convey("A status code of 400 and an error is returned", func() {
			status, err := cli.CreateIndex(context.Background(), testIndex, indexSettings)
			So(err, ShouldNotEqual, nil)
			So(err, ShouldResemble, errorUnexpectedStatusCode)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one")
			So(status, ShouldEqual, 400)
		})
	})
}

func TestDeleteIndex(t *testing.T) {
	Convey("Given that an index is deleted", t, func() {
		httpCli := clientMock(doSuccessful)
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)
		checkClient(httpCli)

		Convey("A status code of 200 and no error is returned ", func() {
			status, err := cli.DeleteIndex(context.Background(), testIndex)
			So(err, ShouldEqual, nil)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one")
			So(status, ShouldEqual, 200)
		})
	})

	Convey("Given that there is a server error", t, func() {

		httpCli := clientMock(doUnsuccessful)
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)
		checkClient(httpCli)

		Convey("A status code of 500 and an error is returned", func() {
			status, err := cli.DeleteIndex(context.Background(), testIndex)
			So(err, ShouldNotEqual, nil)
			So(err, ShouldResemble, ErrUnreachable)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one")
			So(status, ShouldEqual, 0)
		})
	})

	Convey("Given that an elasticsearch returns an unexpected status code", t, func() {

		httpCli := clientMock(unexpectedStatusCode)
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)
		checkClient(httpCli)

		Convey("A status code of 400 and an error is returned", func() {
			status, err := cli.DeleteIndex(context.Background(), testIndex)
			So(err, ShouldNotEqual, nil)
			So(err, ShouldResemble, errorUnexpectedStatusCode)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one")
			So(status, ShouldEqual, 400)
		})
	})
}

func TestAddDocument(t *testing.T) {

	document := []byte("document")

	Convey("Given that an index is created", t, func() {

		httpCli := clientMock(doSuccessful)
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)
		checkClient(httpCli)

		Convey("A status code of 200 and no error is returned", func() {
			status, err := cli.AddDocument(context.Background(), testIndex, testType, testID, document)
			So(err, ShouldEqual, nil)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one/_type/id")
			So(status, ShouldEqual, 200)
		})
	})

	Convey("Given that there is a server error", t, func() {

		httpCli := clientMock(doUnsuccessful)
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)
		checkClient(httpCli)

		Convey("A status code of 500 and an error is returned", func() {
			status, err := cli.AddDocument(context.Background(), testIndex, testType, testID, document)
			So(err, ShouldNotEqual, nil)
			So(err, ShouldResemble, ErrUnreachable)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one/_type/id")
			So(status, ShouldEqual, 0)
		})
	})

	Convey("Given that an elasticsearch returns an unexpected status code", t, func() {

		httpCli := clientMock(unexpectedStatusCode)
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)
		checkClient(httpCli)

		Convey("A status code of 400 and an error is returned", func() {
			status, err := cli.AddDocument(context.Background(), testIndex, testType, testID, document)
			So(err, ShouldNotEqual, nil)
			So(err, ShouldResemble, errorUnexpectedStatusCode)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one/_type/id")
			So(status, ShouldEqual, 400)
		})
	})
}

func checkClient(httpCli *dphttp.ClienterMock) {
	So(httpCli.GetPathsWithNoRetriesCalls(), ShouldHaveLength, 1)
	So(httpCli.SetPathsWithNoRetriesCalls(), ShouldHaveLength, 1)
}

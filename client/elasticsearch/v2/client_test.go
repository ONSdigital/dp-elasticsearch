package v2_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/ONSdigital/dp-elasticsearch/v4/client"
	v2 "github.com/ONSdigital/dp-elasticsearch/v4/client/elasticsearch/v2"
	esError "github.com/ONSdigital/dp-elasticsearch/v4/errors"

	dphttp "github.com/ONSdigital/dp-net/v2/http"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	testType = "_type"
	testID   = "id"

	envAccessKeyID     = "AWS_ACCESS_KEY_ID"
	envSecretAccessKey = "AWS_SECRET_ACCESS_KEY"

	testAccessKey       = "TEST_ACCESS_KEY"
	testSecretAccessKey = "TEST_SECRET_KEY"
)

var (
	ctx = context.Background()

	errorUnexpectedStatusCode = errors.New("unexpected status code from api")

	doSuccessful = func(ctx context.Context, request *http.Request) (*http.Response, error) {
		return resp("do successful", 200), nil
	}

	doSuccessfulCreate = func(ctx context.Context, request *http.Request) (*http.Response, error) {
		return resp("do successful create", 201), nil
	}

	doUnsuccessful = func(ctx context.Context, request *http.Request) (*http.Response, error) {
		return resp("do unsuccessful", 0), ErrUnreachable
	}

	unexpectedStatusCode = func(ctx context.Context, request *http.Request) (*http.Response, error) {
		return resp("unexpected status", 400), nil
	}

	doSuccessfulIndices = func(ctx context.Context, request *http.Request) (*http.Response, error) {
		return resp(`{"ook":"bar"}`, 200), nil
	}

	emptyListOfPathsWithNoRetries = func() []string {
		return []string{}
	}

	setListOfPathsWithNoRetries = func(listOfPaths []string) {}
)

func clientMock(doFunc func(ctx context.Context, request *http.Request) (*http.Response, error)) *dphttp.ClienterMock {
	return &dphttp.ClienterMock{
		DoFunc:                    doFunc,
		GetPathsWithNoRetriesFunc: emptyListOfPathsWithNoRetries,
		SetPathsWithNoRetriesFunc: setListOfPathsWithNoRetries,
	}
}

func TestCreateIndex(t *testing.T) {
	testSetup(t)

	indexSettings := []byte("settings")

	Convey("Given that an index with settings is created", t, func() {
		httpCli := clientMock(doSuccessful)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli)
		checkClient(httpCli)

		Convey("A status code of 200 and no error is returned", func() {
			err := cli.CreateIndex(context.Background(), testIndex, indexSettings)
			So(err, ShouldEqual, nil)
			So(httpCli.DoCalls(), ShouldHaveLength, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one")
		})
	})

	Convey("Given that an index without settings is created", t, func() {
		httpCli := clientMock(doSuccessful)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli)
		checkClient(httpCli)

		Convey("A status code of 200 and no error is returned", func() {
			err := cli.CreateIndex(context.Background(), testIndex, nil)
			So(err, ShouldEqual, nil)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one")
		})
	})

	Convey("Given that there is a server error", t, func() {
		httpCli := clientMock(doUnsuccessful)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli)
		checkClient(httpCli)

		Convey("A status code of 500 and an error is returned", func() {
			err := cli.CreateIndex(context.Background(), testIndex, indexSettings)
			So(err, ShouldNotEqual, nil)
			So(esError.ErrorMessage(err), ShouldEqual, ErrUnreachable.Error())
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one")
		})
	})

	Convey("Given that an elasticsearch returns an unexpected status code", t, func() {
		httpCli := clientMock(unexpectedStatusCode)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli)
		checkClient(httpCli)

		Convey("A status code of 400 and an error is returned", func() {
			err := cli.CreateIndex(context.Background(), testIndex, indexSettings)
			So(err, ShouldNotEqual, nil)
			So(esError.ErrorMessage(err), ShouldEqual, errorUnexpectedStatusCode.Error())
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one")
		})
	})
}

func TestDeleteIndex(t *testing.T) {
	testSetup(t)

	Convey("Given that an index is deleted", t, func() {
		httpCli := clientMock(doSuccessful)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli)
		checkClient(httpCli)

		Convey("A status code of 200 and no error is returned ", func() {
			err := cli.DeleteIndex(context.Background(), testIndex)
			So(err, ShouldEqual, nil)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one")
		})
	})

	Convey("Given that there is a server error", t, func() {
		httpCli := clientMock(doUnsuccessful)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli)
		checkClient(httpCli)

		Convey("A status code of 500 and an error is returned", func() {
			err := cli.DeleteIndex(context.Background(), testIndex)
			So(err, ShouldNotEqual, nil)
			So(esError.ErrorMessage(err), ShouldEqual, ErrUnreachable.Error())
			So(esError.ErrorStatus(err), ShouldEqual, 0)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one")
		})
	})

	Convey("Given that an elasticsearch returns an unexpected status code", t, func() {
		httpCli := clientMock(unexpectedStatusCode)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli)
		checkClient(httpCli)

		Convey("A status code of 400 and an error is returned", func() {
			err := cli.DeleteIndex(context.Background(), testIndex)
			So(err, ShouldNotEqual, nil)
			So(esError.ErrorMessage(err), ShouldEqual, errorUnexpectedStatusCode.Error())
			So(esError.ErrorStatus(err), ShouldEqual, 400)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one")
		})
	})
}

func TestGetIndices(t *testing.T) {
	testSetup(t)
	testIndices := []string{"a", "b"}

	Convey("Given that indices are retrieved", t, func() {
		httpCli := clientMock(doSuccessfulIndices)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli)
		checkClient(httpCli)

		Convey("A status code of 200 and no error is returned ", func() {
			_, err := cli.GetIndices(context.Background(), testIndices)
			So(err, ShouldEqual, nil)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/a,b")
		})
	})

	Convey("Given that there is a server error", t, func() {
		httpCli := clientMock(doUnsuccessful)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli)
		checkClient(httpCli)

		Convey("A status code of 500 and an error is returned", func() {
			_, err := cli.GetIndices(context.Background(), testIndices)
			So(err, ShouldNotEqual, nil)
			So(esError.ErrorMessage(err), ShouldEqual, ErrUnreachable.Error())
			So(esError.ErrorStatus(err), ShouldEqual, 0)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/a,b")
		})
	})
}

func TestAddDocument(t *testing.T) {
	testSetup(t)
	document := []byte("document")

	Convey("Given that an index is created", t, func() {
		httpCli := clientMock(doSuccessfulCreate)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli)
		checkClient(httpCli)

		Convey("No error is returned", func() {
			options := client.AddDocumentOptions{DocumentType: testType}
			err := cli.AddDocument(context.Background(), testIndex, testID, document, &options)
			So(err, ShouldEqual, nil)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one/_type/id")
		})
	})

	Convey("Given that there is a server error", t, func() {
		httpCli := clientMock(doUnsuccessful)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli)
		checkClient(httpCli)

		Convey("An error is returned", func() {
			options := client.AddDocumentOptions{DocumentType: testType}
			err := cli.AddDocument(context.Background(), testIndex, testID, document, &options)
			So(err, ShouldNotEqual, nil)
			So(esError.ErrorMessage(err), ShouldEqual, ErrUnreachable.Error())
			So(esError.ErrorStatus(err), ShouldEqual, 0)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one/_type/id")
		})
	})

	Convey("Given that an elasticsearch returns an unexpected status code", t, func() {
		httpCli := clientMock(unexpectedStatusCode)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli)
		checkClient(httpCli)

		Convey("An error is returned", func() {
			options := client.AddDocumentOptions{DocumentType: testType}
			err := cli.AddDocument(context.Background(), testIndex, testID, document, &options)
			So(err, ShouldNotEqual, nil)
			So(esError.ErrorMessage(err), ShouldEqual, errorUnexpectedStatusCode.Error())
			So(esError.ErrorStatus(err), ShouldEqual, 400)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/one/_type/id")
		})
	})
}

func TestBulkUpdate(t *testing.T) {
	testSetup(t)

	esDestIndex := "ons_test"
	bulk := make([]byte, 1)

	Convey("Given that bulk update is a success", t, func() {
		doFuncWithValidResponse := func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return successESResponse(), nil
		}
		httpCli := clientMock(doFuncWithValidResponse)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli)
		checkClient(httpCli)

		Convey("When bulkupdate is called", func() {
			b, err := cli.BulkUpdate(ctx, esDestIndex, testURL, bulk)

			Convey("Then a status code of 201 and no error is returned ", func() {
				So(err, ShouldEqual, nil)
				So(string(b), ShouldEqual, "Created")
				So(len(httpCli.DoCalls()), ShouldEqual, 1)
				So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/ons_test/_bulk")
			})
		})
	})

	Convey("Given that there is a server error", t, func() {
		doFuncWithInValidResponse := func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return unsuccessfulESResponse(), nil
		}
		httpCli2 := clientMock(doFuncWithInValidResponse)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli2)
		checkClient(httpCli2)

		Convey("When bulkupdate is called", func() {
			_, err := cli.BulkUpdate(ctx, esDestIndex, testURL, bulk)

			Convey("Then a status code of 500 and an error is returned", func() {
				So(err, ShouldNotBeNil)
				So(esError.ErrorMessage(err), ShouldEqual, errors.New("unexpected status code from api").Error())
				So(esError.ErrorStatus(err), ShouldEqual, 500)
				So(len(httpCli2.DoCalls()), ShouldEqual, 1)
				So(httpCli2.DoCalls()[0].Req.URL.Path, ShouldEqual, "/ons_test/_bulk")
			})
		})
	})
}

func checkClient(httpCli *dphttp.ClienterMock) {
	So(httpCli.GetPathsWithNoRetriesCalls(), ShouldHaveLength, 1)
	So(httpCli.SetPathsWithNoRetriesCalls(), ShouldHaveLength, 1)
}

func testSetup(t *testing.T) {
	var err error
	accessKeyID, secretAccessKey := setEnvironmentVars()

	t.Cleanup(func() {
		removeTestEnvironmentVariables(accessKeyID, secretAccessKey)
	})

	if err != nil {
		t.Fatalf("test failed on setup, error: %v", err)
	}
}

func setEnvironmentVars() (accessKeyID, secretAccessKey string) {
	accessKeyID = os.Getenv(envAccessKeyID)
	secretAccessKey = os.Getenv(envSecretAccessKey)

	os.Setenv(envAccessKeyID, testAccessKey)
	os.Setenv(envSecretAccessKey, testSecretAccessKey)

	return
}

func removeTestEnvironmentVariables(accessKeyID, secretAccessKey string) {
	os.Setenv(envAccessKeyID, accessKeyID)
	os.Setenv(envSecretAccessKey, secretAccessKey)
}

func successESResponse() *http.Response {
	return &http.Response{
		StatusCode: 201,
		Body:       io.NopCloser(bytes.NewBufferString(`Created`)),
		Header:     make(http.Header),
	}
}

func unsuccessfulESResponse() *http.Response {
	return &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(bytes.NewBufferString(`Internal server error`)),
		Header:     make(http.Header),
	}
}

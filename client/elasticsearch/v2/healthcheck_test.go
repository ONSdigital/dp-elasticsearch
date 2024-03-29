package v2_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	v2 "github.com/ONSdigital/dp-elasticsearch/v3/client/elasticsearch/v2"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	clusterHealthy       = "{\"cluster_name\" : \"testcluster\", \"status\" : \"green\"}"
	clusterAtRisk        = "{\"cluster_name\" : \"testcluster\", \"status\" : \"yellow\"}"
	clusterUnhealthy     = "{\"cluster_name\" : \"testcluster\", \"status\" : \"red\"}"
	clusterInvalidStatus = "{\"cluster_name\" : \"testcluster\", \"status\" : \"wrongValue\"}"
	clusterMissingStatus = "{\"cluster_name\" : \"testcluster\"}"

	testURL = "http://some.url"

	head = "HEAD"
)

var testIndex = "one"

var testTwoIndexes = "two"

// Error definitions for testing
var (
	ErrUnreachable = errors.New("unreachable")
)

var doOkGreen = func(ctx context.Context, request *http.Request) (*http.Response, error) {
	return resp(clusterHealthy, 200), nil
}

var doOkYellow = func(ctx context.Context, request *http.Request) (*http.Response, error) {
	return resp(clusterAtRisk, 200), nil
}

var doOkRed = func(ctx context.Context, request *http.Request) (*http.Response, error) {
	return resp(clusterUnhealthy, 200), nil
}

var doOkInvalidStatus = func(ctx context.Context, request *http.Request) (*http.Response, error) {
	return resp(clusterInvalidStatus, 200), nil
}

var doOkMissingStatus = func(ctx context.Context, request *http.Request) (*http.Response, error) {
	return resp(clusterMissingStatus, 200), nil
}

var doUnexpectedCode = func(ctx context.Context, request *http.Request) (*http.Response, error) {
	return resp(clusterHealthy, 300), nil
}

var doUnreachable = func(ctx context.Context, request *http.Request) (*http.Response, error) {
	return nil, ErrUnreachable
}

var doIndexExists = func(ctx context.Context, request *http.Request) (*http.Response, error) {
	return indexResp(200), nil
}

var doIndexDoesNotExist = func(ctx context.Context, request *http.Request) (*http.Response, error) {
	return indexResp(404), nil
}

var doUnexpectedIndexResponse = func(ctx context.Context, request *http.Request) (*http.Response, error) {
	return indexResp(300), nil
}

func resp(body string, code int) *http.Response {
	return &http.Response{
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		StatusCode: code,
	}
}

func indexResp(code int) *http.Response {
	return &http.Response{
		Body:       io.NopCloser(bytes.NewBuffer([]byte{})),
		StatusCode: code,
	}
}

func TestElasticsearchHealthGreen(t *testing.T) {
	testSetup(t)
	Convey("Given that Elasticsearch is healthy", t, func() {
		httpCli := clientMock(
			func(ctx context.Context, request *http.Request) (*http.Response, error) {
				if request.Method == head {
					return doIndexExists(ctx, request)
				}
				return doOkGreen(ctx, request)
			},
		)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli, testIndex)
		checkClient(httpCli)

		// CheckState for test validation
		checkState := health.NewCheckState(v2.ServiceName)

		Convey("Checker updates the CheckState to a successful state", func() {
			cli.Checker(context.Background(), checkState)
			So(len(httpCli.DoCalls()), ShouldEqual, 2)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/_cluster/health")
			So(httpCli.DoCalls()[1].Req.URL.Path, ShouldEqual, "/one")
			So(checkState.Status(), ShouldEqual, health.StatusOK)
			So(checkState.Message(), ShouldEqual, v2.MsgHealthy)
			So(checkState.StatusCode(), ShouldEqual, 200)
		})
	})
}

func TestElasticsearchHealthYellow(t *testing.T) {
	testSetup(t)
	Convey("Given that Elasticsearch data is at risk", t, func() {
		httpCli := clientMock(doOkYellow)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli, testIndex)
		checkClient(httpCli)

		// CheckState for test validation
		checkState := health.NewCheckState(v2.ServiceName)

		Convey("Checker updates the CheckState to a Warning state structure with the relevant error message", func() {
			cli.Checker(context.Background(), checkState)
			So(len(httpCli.DoCalls()), ShouldEqual, 2)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/_cluster/health")
			So(checkState.Status(), ShouldEqual, health.StatusOK)
			So(checkState.Message(), ShouldEqual, v2.ErrorClusterAtRisk.Error())
			So(checkState.StatusCode(), ShouldEqual, 200)
		})
	})
}

func TestElasticsearchHealthRed(t *testing.T) {
	testSetup(t)
	Convey("Given that Elasticsearch is unhealthy", t, func() {
		httpCli := clientMock(doOkRed)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli, testIndex)
		checkClient(httpCli)

		// CheckState for test validation
		checkState := health.NewCheckState(v2.ServiceName)

		Convey("Checker updates the CheckState to a Critical state with the relevant error message", func() {
			cli.Checker(context.Background(), checkState)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/_cluster/health")
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, v2.ErrorUnhealthyClusterStatus.Error())
			So(checkState.StatusCode(), ShouldEqual, 200)
		})
	})
}

func TestElasticsearchInvalidHealth(t *testing.T) {
	testSetup(t)
	Convey("Given that Elasticsearch API returns an invalid status", t, func() {
		httpCli := clientMock(doOkInvalidStatus)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli, testIndex)
		checkClient(httpCli)

		// CheckState for test validation
		checkState := health.NewCheckState(v2.ServiceName)

		Convey("Checker updates the CheckState to a critical state with the relevant error message", func() {
			cli.Checker(context.Background(), checkState)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/_cluster/health")
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, v2.ErrorInvalidHealthStatus.Error())
			So(checkState.StatusCode(), ShouldEqual, 200)
		})
	})
}

func TestElasticsearchMissingHealth(t *testing.T) {
	testSetup(t)
	Convey("Given that Elasticsearch API response does not provide the health status", t, func() {
		httpCli := clientMock(doOkMissingStatus)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli, testIndex)
		checkClient(httpCli)

		// CheckState for test validation
		checkState := health.NewCheckState(v2.ServiceName)

		Convey("Checker updates the CheckState to a critical state with the relevant error message", func() {
			cli.Checker(context.Background(), checkState)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/_cluster/health")
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, v2.ErrorInvalidHealthStatus.Error())
			So(checkState.StatusCode(), ShouldEqual, 200)
		})
	})
}

func TestUnexpectedStatusCode(t *testing.T) {
	testSetup(t)
	Convey("Given that Elasticsearch API response provides a wrong Status Code", t, func() {
		httpCli := clientMock(doUnexpectedCode)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli, testIndex)
		checkClient(httpCli)

		// CheckState for test validation
		checkState := health.NewCheckState(v2.ServiceName)

		Convey("Checker updates the CheckState to a critical state with the relevant error message", func() {
			cli.Checker(context.Background(), checkState)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/_cluster/health")
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, v2.ErrorUnexpectedStatusCode.Error())
			So(checkState.StatusCode(), ShouldEqual, 300)
		})
	})
}

func TestExceptionUnreachable(t *testing.T) {
	testSetup(t)
	Convey("Given that Elasticsearch is unreachable", t, func() {
		httpCli := clientMock(doUnreachable)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli, testIndex)
		checkClient(httpCli)

		// CheckState for test validation
		checkState := health.NewCheckState(v2.ServiceName)

		Convey("Checker updates the CheckState to a critical state with the relevant error message", func() {
			cli.Checker(context.Background(), checkState)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/_cluster/health")
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, ErrUnreachable.Error())
			So(checkState.StatusCode(), ShouldEqual, 500)
		})
	})
}

func TestIndexExists(t *testing.T) {
	testSetup(t)
	Convey("Given that the client has one index and this index exists", t, func() {
		httpCli := clientMock(
			func(ctx context.Context, request *http.Request) (*http.Response, error) {
				if request.Method == head {
					return doIndexExists(ctx, request)
				}
				return doOkGreen(ctx, request)
			},
		)

		cli := v2.NewClientWithHTTPClient(testURL, httpCli, testIndex)
		checkClient(httpCli)

		// CheckState for test validation
		checkState := health.NewCheckState(v2.ServiceName)

		Convey("Checker updates the CheckState to a successful state", func() {
			cli.Checker(context.Background(), checkState)
			So(len(httpCli.DoCalls()), ShouldEqual, 2)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/_cluster/health")
			So(httpCli.DoCalls()[1].Req.URL.Path, ShouldEqual, "/one")
			So(checkState.Status(), ShouldEqual, health.StatusOK)
			So(checkState.Message(), ShouldEqual, v2.MsgHealthy)
			So(checkState.StatusCode(), ShouldEqual, 200)
		})
	})
}

func TestIndexDoesNotExist(t *testing.T) {
	testSetup(t)
	Convey("Given that the client has one index and this index does not exists", t, func() {
		httpCli := clientMock(
			func(ctx context.Context, request *http.Request) (*http.Response, error) {
				if request.Method == head {
					return doIndexDoesNotExist(ctx, request)
				}
				return doOkGreen(ctx, request)
			},
		)

		cli := v2.NewClientWithHTTPClient(testURL, httpCli, testIndex)
		checkClient(httpCli)

		// CheckState for test validation
		checkState := health.NewCheckState(v2.ServiceName)

		Convey("Checker updates the CheckState to a successful state", func() {
			cli.Checker(context.Background(), checkState)
			So(len(httpCli.DoCalls()), ShouldEqual, 2)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/_cluster/health")
			So(httpCli.DoCalls()[1].Req.URL.Path, ShouldEqual, "/one")
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, v2.ErrorIndexDoesNotExist.Error())
			So(checkState.StatusCode(), ShouldEqual, 404)
		})
	})
}

func TestNoClientIndex(t *testing.T) {
	testSetup(t)
	Convey("Given that the client does not have any indexes", t, func() {
		httpCli := clientMock(
			func(ctx context.Context, request *http.Request) (*http.Response, error) {
				return doOkGreen(ctx, request)
			},
		)

		cli := v2.NewClientWithHTTPClient(testURL, httpCli)
		checkClient(httpCli)

		// CheckState for test validation
		checkState := health.NewCheckState(v2.ServiceName)

		Convey("Checker updates the CheckState to a successful state", func() {
			cli.Checker(context.Background(), checkState)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/_cluster/health")
			So(checkState.Status(), ShouldEqual, health.StatusOK)
			So(checkState.Message(), ShouldEqual, v2.MsgHealthy)
			So(checkState.StatusCode(), ShouldEqual, 200)
		})
	})
}

func TestTwoIndexesExist(t *testing.T) {
	testSetup(t)

	Convey("Given that the client has two indexes and both indexes exist", t, func() {
		httpCli := clientMock(
			func(ctx context.Context, request *http.Request) (*http.Response, error) {
				if request.Method == head {
					return doIndexExists(ctx, request)
				}
				return doOkGreen(ctx, request)
			},
		)

		cli := v2.NewClientWithHTTPClient(testURL, httpCli, testIndex, testTwoIndexes)
		checkClient(httpCli)

		// CheckState for test validation
		checkState := health.NewCheckState(v2.ServiceName)

		Convey("Checker updates the CheckState to a successful state", func() {
			cli.Checker(context.Background(), checkState)
			So(len(httpCli.DoCalls()), ShouldEqual, 3)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/_cluster/health")
			So(httpCli.DoCalls()[1].Req.URL.Path, ShouldEqual, "/one")
			So(httpCli.DoCalls()[2].Req.URL.Path, ShouldEqual, "/two")
			So(checkState.Status(), ShouldEqual, health.StatusOK)
			So(checkState.Message(), ShouldEqual, v2.MsgHealthy)
			So(checkState.StatusCode(), ShouldEqual, 200)
		})
	})
}

func TestOneOfTwoIndexesExist(t *testing.T) {
	testSetup(t)

	Convey("Given that the client has two indexes and only the first index exists", t, func() {
		httpCli := clientMock(
			func(ctx context.Context, request *http.Request) (*http.Response, error) {
				if request.Method == head {
					if request.URL.Path == "/one" {
						return doIndexExists(ctx, request)
					}
					return doIndexDoesNotExist(ctx, request)
				}
				return doOkGreen(ctx, request)
			},
		)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli, testIndex, testTwoIndexes)
		checkClient(httpCli)

		// CheckState for test validation
		checkState := health.NewCheckState(v2.ServiceName)

		Convey("Checker updates the CheckState to a successful state", func() {
			cli.Checker(context.Background(), checkState)
			So(len(httpCli.DoCalls()), ShouldEqual, 3)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/_cluster/health")
			So(httpCli.DoCalls()[1].Req.URL.Path, ShouldEqual, "/one")
			So(httpCli.DoCalls()[2].Req.URL.Path, ShouldEqual, "/two")
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, v2.ErrorIndexDoesNotExist.Error())
			So(checkState.StatusCode(), ShouldEqual, 404)
		})
	})
}

func TestUnexpectedIndexResponse(t *testing.T) {
	testSetup(t)

	Convey("Given that the elasticsearch response is unexpected", t, func() {
		httpCli := clientMock(
			func(ctx context.Context, request *http.Request) (*http.Response, error) {
				if request.Method == head {
					return doUnexpectedIndexResponse(ctx, request)
				}
				return doOkGreen(ctx, request)
			},
		)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli, testIndex)
		checkClient(httpCli)

		// CheckState for test validation
		checkState := health.NewCheckState(v2.ServiceName)

		Convey("Checker updates the CheckState to a successful state", func() {
			cli.Checker(context.Background(), checkState)
			So(len(httpCli.DoCalls()), ShouldEqual, 2)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/_cluster/health")
			So(httpCli.DoCalls()[1].Req.URL.Path, ShouldEqual, "/one")
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, v2.ErrorUnexpectedStatusCode.Error())
			So(checkState.StatusCode(), ShouldEqual, 300)
		})
	})
}

func TestSecondExceptionUnreachable(t *testing.T) {
	testSetup(t)

	Convey("Given that elasticsearch is unreachable", t, func() {
		httpCli := clientMock(
			func(ctx context.Context, request *http.Request) (*http.Response, error) {
				if request.Method == head {
					return doUnreachable(ctx, request)
				}
				return doOkGreen(ctx, request)
			},
		)
		cli := v2.NewClientWithHTTPClient(testURL, httpCli, testIndex)
		checkClient(httpCli)

		// CheckState for test validation
		checkState := health.NewCheckState(v2.ServiceName)

		Convey("Checker updates the CheckState to a successful state", func() {
			cli.Checker(context.Background(), checkState)
			So(len(httpCli.DoCalls()), ShouldEqual, 2)
			So(httpCli.DoCalls()[0].Req.URL.Path, ShouldEqual, "/_cluster/health")
			So(httpCli.DoCalls()[1].Req.URL.Path, ShouldEqual, "/one")
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, ErrUnreachable.Error())
			So(checkState.StatusCode(), ShouldEqual, 500)
		})
	})
}

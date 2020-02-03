package elasticsearch_test

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	elasticsearch "github.com/ONSdigital/dp-elasticsearch"
	"github.com/ONSdigital/dp-elasticsearch/mock"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	clusterHealthy       = "{\"cluster_name\" : \"testcluster\", \"status\" : \"green\"}"
	clusterAtRisk        = "{\"cluster_name\" : \"testcluster\", \"status\" : \"yellow\"}"
	clusterUnhealthy     = "{\"cluster_name\" : \"testcluster\", \"status\" : \"red\"}"
	clusterInvalidStatus = "{\"cluster_name\" : \"testcluster\", \"status\" : \"wrongValue\"}"
	clusterMissingStatus = "{\"cluster_name\" : \"testcluster\"}"
)

const testUrl = "http://some.url"

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

func resp(body string, code int) *http.Response {
	return &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		StatusCode: code,
	}
}

func TestElasticsearchHealthGreen(t *testing.T) {
	Convey("Given that Elasticsearch is healthy", t, func() {

		var httpCli = &mock.RchttpClientMock{
			DoFunc: doOkGreen,
		}
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)

		// mock CheckState for test validation
		mockCheckState := mock.CheckStateMock{
			UpdateFunc: func(status, message string, statusCode int) error {
				return nil
			},
		}

		Convey("Checker updates the CheckState to a successful state", func() {
			cli.Checker(context.Background(), &mockCheckState)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			updateCalls := mockCheckState.UpdateCalls()
			So(len(updateCalls), ShouldEqual, 1)
			So(updateCalls[0].Status, ShouldEqual, health.StatusOK)
			So(updateCalls[0].Message, ShouldEqual, elasticsearch.MsgHealthy)
			So(updateCalls[0].StatusCode, ShouldEqual, 200)
		})
	})
}

func TestElasticsearchHealthYellow(t *testing.T) {
	Convey("Given that Elasticsearch data is at risk", t, func() {

		var httpCli = &mock.RchttpClientMock{
			DoFunc: doOkYellow,
		}
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)

		// mock CheckState for test validation
		mockCheckState := mock.CheckStateMock{
			UpdateFunc: func(status, message string, statusCode int) error {
				return nil
			},
		}

		Convey("Checker updates the CheckState to a Warning state structure with the relevant error message", func() {
			cli.Checker(context.Background(), &mockCheckState)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			updateCalls := mockCheckState.UpdateCalls()
			So(len(updateCalls), ShouldEqual, 1)
			So(updateCalls[0].Status, ShouldEqual, health.StatusWarning)
			So(updateCalls[0].Message, ShouldEqual, elasticsearch.ErrorClusterAtRisk.Error())
			So(updateCalls[0].StatusCode, ShouldEqual, 200)
		})
	})
}

func TestElasticsearchHealthRed(t *testing.T) {
	Convey("Given that Elasticsearch is unhealthy", t, func() {

		var httpCli = &mock.RchttpClientMock{
			DoFunc: doOkRed,
		}
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)

		// mock CheckState for test validation
		mockCheckState := mock.CheckStateMock{
			UpdateFunc: func(status, message string, statusCode int) error {
				return nil
			},
		}

		Convey("Checker updates the CheckState to a Critical state with the relevant error message", func() {
			cli.Checker(context.Background(), &mockCheckState)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			updateCalls := mockCheckState.UpdateCalls()
			So(len(updateCalls), ShouldEqual, 1)
			So(updateCalls[0].Status, ShouldEqual, health.StatusCritical)
			So(updateCalls[0].Message, ShouldEqual, elasticsearch.ErrorUnhealthyClusterStatus.Error())
			So(updateCalls[0].StatusCode, ShouldEqual, 200)
		})
	})
}

func TestElasticsearchInvalidHealth(t *testing.T) {
	Convey("Given that Elasticsearch API returns an invalid status", t, func() {

		var httpCli = &mock.RchttpClientMock{
			DoFunc: doOkInvalidStatus,
		}
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)

		// mock CheckState for test validation
		mockCheckState := mock.CheckStateMock{
			UpdateFunc: func(status, message string, statusCode int) error {
				return nil
			},
		}

		Convey("Checker updates the CheckState to a critical state with the relevant error message", func() {
			cli.Checker(context.Background(), &mockCheckState)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			updateCalls := mockCheckState.UpdateCalls()
			So(len(updateCalls), ShouldEqual, 1)
			So(updateCalls[0].Status, ShouldEqual, health.StatusCritical)
			So(updateCalls[0].Message, ShouldEqual, elasticsearch.ErrorInvalidHealthStatus.Error())
			So(updateCalls[0].StatusCode, ShouldEqual, 200)
		})
	})
}

func TestElasticsearchMissingHealth(t *testing.T) {
	Convey("Given that Elasticsearch API response does not provide the health status", t, func() {

		var httpCli = &mock.RchttpClientMock{
			DoFunc: doOkMissingStatus,
		}
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)

		// mock CheckState for test validation
		mockCheckState := mock.CheckStateMock{
			UpdateFunc: func(status, message string, statusCode int) error {
				return nil
			},
		}

		Convey("Checker updates the CheckState to a critical state with the relevant error message", func() {
			cli.Checker(context.Background(), &mockCheckState)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			updateCalls := mockCheckState.UpdateCalls()
			So(len(updateCalls), ShouldEqual, 1)
			So(updateCalls[0].Status, ShouldEqual, health.StatusCritical)
			So(updateCalls[0].Message, ShouldEqual, elasticsearch.ErrorInvalidHealthStatus.Error())
			So(updateCalls[0].StatusCode, ShouldEqual, 200)
		})
	})
}

func TestUnexpectedStatusCode(t *testing.T) {
	Convey("Given that Elasticsearch API response provides a wrong Status Code", t, func() {

		var httpCli = &mock.RchttpClientMock{
			DoFunc: doUnexpectedCode,
		}
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)

		// mock CheckState for test validation
		mockCheckState := mock.CheckStateMock{
			UpdateFunc: func(status, message string, statusCode int) error {
				return nil
			},
		}

		Convey("Checker updates the CheckState to a critical state with the relevant error message", func() {
			cli.Checker(context.Background(), &mockCheckState)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			updateCalls := mockCheckState.UpdateCalls()
			So(len(updateCalls), ShouldEqual, 1)
			So(updateCalls[0].Status, ShouldEqual, health.StatusCritical)
			So(updateCalls[0].Message, ShouldEqual, elasticsearch.ErrorUnexpectedStatusCode.Error())
			So(updateCalls[0].StatusCode, ShouldEqual, 300)
		})
	})
}

func TestExceptionUnreachable(t *testing.T) {
	Convey("Given that Elasticsearch is unreachable", t, func() {

		var httpCli = &mock.RchttpClientMock{
			DoFunc: doUnreachable,
		}
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)

		// mock CheckState for test validation
		mockCheckState := mock.CheckStateMock{
			UpdateFunc: func(status, message string, statusCode int) error {
				return nil
			},
		}

		Convey("Checker updates the CheckState to a critical state with the relevant error message", func() {
			cli.Checker(context.Background(), &mockCheckState)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
			updateCalls := mockCheckState.UpdateCalls()
			So(len(updateCalls), ShouldEqual, 1)
			So(updateCalls[0].Status, ShouldEqual, health.StatusCritical)
			So(updateCalls[0].Message, ShouldEqual, ErrUnreachable.Error())
			So(updateCalls[0].StatusCode, ShouldEqual, 500)
		})
	})
}

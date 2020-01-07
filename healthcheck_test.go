package elasticsearch_test

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

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
	ErrUnreacheable = errors.New("unreacheable")
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

var doUnreacheable = func(ctx context.Context, request *http.Request) (*http.Response, error) {
	return nil, ErrUnreacheable
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

		Convey("Checker returns a successful Check structure", func() {
			validateSuccessfulCheck(cli)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
		})
	})
}

func TestElasticsearchHealthYellow(t *testing.T) {
	Convey("Given that Elasticsearch data is at risk", t, func() {

		var httpCli = &mock.RchttpClientMock{
			DoFunc: doOkYellow,
		}
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)

		Convey("Checker returns a Warning Check structure", func() {
			_, err := validateWarningCheck(cli, 200, elasticsearch.ErrorClusterAtRisk.Error())
			So(err, ShouldEqual, elasticsearch.ErrorClusterAtRisk)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
		})
	})
}

func TestElasticsearchHealthRed(t *testing.T) {
	Convey("Given that Elasticsearch is unhealthy", t, func() {

		var httpCli = &mock.RchttpClientMock{
			DoFunc: doOkRed,
		}
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)

		Convey("Checker returns a Critical Check structure", func() {
			_, err := validateCriticalCheck(cli, 200, elasticsearch.ErrorUnhealthyClusterStatus.Error())
			So(err, ShouldEqual, elasticsearch.ErrorUnhealthyClusterStatus)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
		})
	})
}

func TestElasticsearchInvalidHealth(t *testing.T) {
	Convey("Given that Elasticsearch API returns an invalid status", t, func() {

		var httpCli = &mock.RchttpClientMock{
			DoFunc: doOkInvalidStatus,
		}
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)

		Convey("Checker returns a critical Check structure", func() {
			_, err := validateCriticalCheck(cli, 200, elasticsearch.ErrorInvalidHealthStatus.Error())
			So(err, ShouldEqual, elasticsearch.ErrorInvalidHealthStatus)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
		})
	})
}

func TestElasticsearchMissingHealth(t *testing.T) {
	Convey("Given that Elasticsearch API response does not provide the health status", t, func() {

		var httpCli = &mock.RchttpClientMock{
			DoFunc: doOkMissingStatus,
		}
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)

		Convey("Checker returns a critical Check structure", func() {
			_, err := validateCriticalCheck(cli, 200, elasticsearch.ErrorInvalidHealthStatus.Error())
			So(err, ShouldEqual, elasticsearch.ErrorInvalidHealthStatus)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
		})
	})
}

func TestUnexpectedStatusCode(t *testing.T) {
	Convey("Given that Elasticsearch API response provides a wrong Status Code", t, func() {

		var httpCli = &mock.RchttpClientMock{
			DoFunc: doUnexpectedCode,
		}
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)

		Convey("Checker returns a critical Check structure", func() {
			_, err := validateCriticalCheck(cli, 300, elasticsearch.ErrorUnexpectedStatusCode.Error())
			So(err, ShouldEqual, elasticsearch.ErrorUnexpectedStatusCode)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
		})
	})
}

func TestExceptionUnreacheable(t *testing.T) {
	Convey("Given that Elasticsearch is unreacheable", t, func() {

		var httpCli = &mock.RchttpClientMock{
			DoFunc: doUnreacheable,
		}
		cli := elasticsearch.NewClientWithHTTPClient(testUrl, true, httpCli)

		Convey("Checker returns a critical Check structure", func() {
			_, err := validateCriticalCheck(cli, 500, ErrUnreacheable.Error())
			So(err, ShouldNotBeNil)
			So(len(httpCli.DoCalls()), ShouldEqual, 1)
		})
	})
}

func validateSuccessfulCheck(cli *elasticsearch.Client) (check *health.Check) {
	ctx := context.Background()
	t0 := time.Now().UTC()
	check, err := cli.Checker(&ctx)
	t1 := time.Now().UTC()
	So(err, ShouldBeNil)
	So(check.Name, ShouldEqual, elasticsearch.ServiceName)
	So(check.Status, ShouldEqual, health.StatusOK)
	So(check.StatusCode, ShouldEqual, 200)
	So(check.Message, ShouldEqual, elasticsearch.MsgHealthy)
	So(check.LastChecked, ShouldHappenOnOrBetween, t0, t1)
	So(check.LastSuccess, ShouldHappenOnOrBetween, t0, t1)
	So(check.LastFailure, ShouldHappenBefore, t0)
	return check
}

func validateWarningCheck(cli *elasticsearch.Client, expectedCode int, expectedMessage string) (check *health.Check, err error) {
	ctx := context.Background()
	t0 := time.Now().UTC()
	check, err = cli.Checker(&ctx)
	t1 := time.Now().UTC()
	So(check.Name, ShouldEqual, elasticsearch.ServiceName)
	So(check.Status, ShouldEqual, health.StatusWarning)
	So(check.StatusCode, ShouldEqual, expectedCode)
	So(check.Message, ShouldEqual, expectedMessage)
	So(check.LastChecked, ShouldHappenOnOrBetween, t0, t1)
	So(check.LastSuccess, ShouldHappenBefore, t0)
	So(check.LastFailure, ShouldHappenOnOrBetween, t0, t1)
	return check, err
}

func validateCriticalCheck(cli *elasticsearch.Client, expectedCode int, expectedMessage string) (check *health.Check, err error) {
	ctx := context.Background()
	t0 := time.Now().UTC()
	check, err = cli.Checker(&ctx)
	t1 := time.Now().UTC()
	So(check.Name, ShouldEqual, elasticsearch.ServiceName)
	So(check.Status, ShouldEqual, health.StatusCritical)
	So(check.StatusCode, ShouldEqual, expectedCode)
	So(check.Message, ShouldEqual, expectedMessage)
	So(check.LastChecked, ShouldHappenOnOrBetween, t0, t1)
	So(check.LastSuccess, ShouldHappenBefore, t0)
	So(check.LastFailure, ShouldHappenOnOrBetween, t0, t1)
	return check, err
}

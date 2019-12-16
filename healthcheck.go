package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"

	rchttp "github.com/ONSdigital/dp-rchttp"
	"github.com/ONSdigital/log.go/log"

	"net/http"
	"net/url"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	awsauth "github.com/smartystreets/go-aws-auth"
)

// List of errors
var (
	ErrorUnexpectedStatusCode   = errors.New("unexpected status code from api")
	ErrorParsingBody            = errors.New("error parsing cluster health response body")
	ErrorUnhealthyClusterStatus = errors.New("error cluster health red")
)

// StatusDescription : Map of descriptions by status
var StatusDescription = map[string]string{
	health.StatusOK:       "Everything is ok",
	health.StatusWarning:  "Things are degraded, but at least partially functioning",
	health.StatusCritical: "The checked functionality is unavailable or non-functioning",
}

// minTime : Oldest time for Check structure.
var minTime = time.Unix(0, 0)

const unhealthy = "red"

// HealthCheckClient provides a healthcheck.Client implementation for health checking elasticsearch.
type HealthCheckClient struct {
	cli          *rchttp.Client
	path         string
	serviceName  string
	signRequests bool
	check        *health.Check
}

// ClusterHealth represents the response from the elasticsearch cluster health check
type ClusterHealth struct {
	Status string `json:"status"`
}

// NewHealthCheckClient returns a new elasticsearch health check client.
func NewHealthCheckClient(url string, signRequests bool) *HealthCheckClient {

	return &HealthCheckClient{
		cli:          rchttp.DefaultClient,
		path:         url + "/_cluster/health",
		serviceName:  "elasticsearch",
		signRequests: signRequests,
	}
}

// healthcheck calls elasticsearch to check its health status.
func (elasticsearch *HealthCheckClient) healthcheck() (string, error) {

	logData := log.Data{"url": elasticsearch.path}
	log.Event(nil, "Created vault client", logData)

	URL, err := url.Parse(elasticsearch.path)
	if err != nil {
		log.Event(nil, "failed to create url for elasticsearch healthcheck", logData, log.Error(err))
		return elasticsearch.serviceName, err
	}

	path := URL.String()
	logData["url"] = path

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Event(nil, "failed to create request for healthcheck call to elasticsearch", logData, log.Error(err))
		return elasticsearch.serviceName, err
	}

	if elasticsearch.signRequests {
		awsauth.Sign(req)
	}

	resp, err := elasticsearch.cli.Do(context.Background(), req)
	if err != nil {
		log.Event(nil, "failed to call elasticsearch", logData, log.Error(err))
		return elasticsearch.serviceName, err
	}
	defer resp.Body.Close()

	logData["http_code"] = resp.StatusCode
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		log.Event(nil, "", logData, log.Error(ErrorUnexpectedStatusCode))
		return elasticsearch.serviceName, ErrorUnexpectedStatusCode
	}

	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Event(nil, "failed to read response body from call to elastic", logData, log.Error(err))
		return elasticsearch.serviceName, ErrorUnexpectedStatusCode
	}

	var clusterHealth ClusterHealth
	err = json.Unmarshal(jsonBody, &clusterHealth)
	if err != nil {
		log.Event(nil, "", logData, log.Error(ErrorParsingBody))
		return elasticsearch.serviceName, ErrorParsingBody
	}

	logData["cluster_health"] = clusterHealth.Status
	if clusterHealth.Status == unhealthy {
		log.Event(nil, "", logData, log.Error(ErrorUnhealthyClusterStatus))
		return elasticsearch.serviceName, ErrorUnhealthyClusterStatus
	}

	return elasticsearch.serviceName, nil
}

// Checker : Check health of Elasticsearch and return it inside a Check structure
func (elasticsearch *HealthCheckClient) Checker(ctx *context.Context) (*health.Check, error) {
	_, err := elasticsearch.healthcheck()
	if err != nil {
		switch err {
		case ErrorUnexpectedStatusCode:
			elasticsearch.check = elasticsearch.getCheck(ctx, 429)
		default:
			elasticsearch.check = elasticsearch.getCheck(ctx, 500)
		}
		return elasticsearch.check, err
	}
	elasticsearch.check = elasticsearch.getCheck(ctx, 200)
	return elasticsearch.check, nil
}

// getCheck : Create a Check structure and populate it according to the error
func (elasticsearch *HealthCheckClient) getCheck(ctx *context.Context, code int) *health.Check {

	currentTime := time.Now().UTC()

	check := &health.Check{
		Name:        elasticsearch.serviceName,
		StatusCode:  code,
		LastChecked: currentTime,
		LastSuccess: minTime,
		LastFailure: minTime,
	}

	switch code {
	case 200:
		check.Message = StatusDescription[health.StatusOK]
		check.Status = health.StatusOK
		check.LastSuccess = currentTime
	case 429:
		check.Message = StatusDescription[health.StatusWarning]
		check.Status = health.StatusWarning
		check.LastFailure = currentTime
	default:
		check.Message = StatusDescription[health.StatusCritical]
		check.Status = health.StatusCritical
		check.LastFailure = currentTime
	}

	return check
}

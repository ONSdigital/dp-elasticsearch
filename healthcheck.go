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

// Client is an ElasticSearch client with a health.Check structure implementation for health checking elasticsearch.
type Client struct {
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

// NewHealthCheckClient returns a new initialized elasticsearch client, without health.Check.
func NewHealthCheckClient(url string, signRequests bool) *Client {

	return &Client{
		cli:          rchttp.DefaultClient,
		path:         url + "/_cluster/health",
		serviceName:  "elasticsearch",
		signRequests: signRequests,
	}
}

// healthcheck calls elasticsearch to check its health status. This call implements only the logic,
// without providing the Check object, and it's aimed for internal use.
func (cli *Client) healthcheck() (string, error) {

	logData := log.Data{"url": cli.path}
	log.Event(nil, "Created vault client", logData)

	URL, err := url.Parse(cli.path)
	if err != nil {
		log.Event(nil, "failed to create url for elasticsearch healthcheck", logData, log.Error(err))
		return cli.serviceName, err
	}

	path := URL.String()
	logData["url"] = path

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Event(nil, "failed to create request for healthcheck call to elasticsearch", logData, log.Error(err))
		return cli.serviceName, err
	}

	if cli.signRequests {
		awsauth.Sign(req)
	}

	resp, err := cli.cli.Do(context.Background(), req)
	if err != nil {
		log.Event(nil, "failed to call elasticsearch", logData, log.Error(err))
		return cli.serviceName, err
	}
	defer resp.Body.Close()

	logData["http_code"] = resp.StatusCode
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		log.Event(nil, "", logData, log.Error(ErrorUnexpectedStatusCode))
		return cli.serviceName, ErrorUnexpectedStatusCode
	}

	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Event(nil, "failed to read response body from call to elastic", logData, log.Error(err))
		return cli.serviceName, ErrorUnexpectedStatusCode
	}

	var clusterHealth ClusterHealth
	err = json.Unmarshal(jsonBody, &clusterHealth)
	if err != nil {
		log.Event(nil, "", logData, log.Error(ErrorParsingBody))
		return cli.serviceName, ErrorParsingBody
	}

	logData["cluster_health"] = clusterHealth.Status
	if clusterHealth.Status == unhealthy {
		log.Event(nil, "", logData, log.Error(ErrorUnhealthyClusterStatus))
		return cli.serviceName, ErrorUnhealthyClusterStatus
	}

	return cli.serviceName, nil
}

// Checker : Check health of Elasticsearch and return it inside a Check structure
func (cli *Client) Checker(ctx *context.Context) (*health.Check, error) {
	_, err := cli.healthcheck()
	if err != nil {
		switch err {
		case ErrorUnexpectedStatusCode:
			cli.check = cli.getCheck(ctx, 429)
		default:
			cli.check = cli.getCheck(ctx, 500)
		}
		return cli.check, err
	}
	cli.check = cli.getCheck(ctx, 200)
	return cli.check, nil
}

// getCheck : Create a Check structure and populate it according to the error
func (cli *Client) getCheck(ctx *context.Context, code int) *health.Check {

	currentTime := time.Now().UTC()

	check := &health.Check{
		Name:        cli.serviceName,
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

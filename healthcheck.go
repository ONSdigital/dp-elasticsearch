package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/log"
	awsauth "github.com/smartystreets/go-aws-auth"
)

// HTTP path to check the health of a cluster using Elasticsearch API
const pathHealth = "/_cluster/health"

// HealthStatus - iota enum of possible health states returned by Elasticsearch API
type HealthStatus int

// Possible values for the HealthStatus
const (
	HealthGreen = iota
	HealthYellow
	HealthRed
)

var healthValues = []string{"green", "yellow", "red"}

func (hs HealthStatus) String() string {
	return healthValues[hs]
}

// List of errors
var (
	ErrorUnexpectedStatusCode   = errors.New("unexpected status code from api")
	ErrorParsingBody            = errors.New("error parsing cluster health response body")
	ErrorClusterAtRisk          = errors.New("error cluster state yellow. Data might be at risk, check your replica shards")
	ErrorUnhealthyClusterStatus = errors.New("error cluster health red. Cluster is unhealthy")
	ErrorInvalidHealthStatus    = errors.New("error invalid health status returned")
)

// StatusDescription : Map of descriptions by status
var StatusDescription = map[string]string{
	health.StatusOK:       "Everything is ok",
	health.StatusWarning:  "Things are degraded, but at least partially functioning",
	health.StatusCritical: "The checked functionality is unavailable or non-functioning",
}

// minTime : Oldest time for Check structure.
var minTime = time.Unix(0, 0)

// ClusterHealth represents the response from the elasticsearch cluster health check
type ClusterHealth struct {
	Status string `json:"status"`
}

// healthcheck calls elasticsearch to check its health status. This call implements only the logic,
// without providing the Check object, and it's aimed for internal use.
func (cli *Client) healthcheck() (string, error) {

	urlHealth := cli.url + pathHealth
	logData := log.Data{"url": urlHealth}

	URL, err := url.Parse(urlHealth)
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

	resp, err := cli.httpCli.Do(context.Background(), req)
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
	switch clusterHealth.Status {
	case healthValues[HealthGreen]:
		return cli.serviceName, nil
	case healthValues[HealthYellow]:
		log.Event(nil, "", logData, log.Error(ErrorClusterAtRisk))
		return cli.serviceName, ErrorClusterAtRisk
	case healthValues[HealthRed]:
		log.Event(nil, "", logData, log.Error(ErrorUnhealthyClusterStatus))
		return cli.serviceName, ErrorUnhealthyClusterStatus
	default:
		log.Event(nil, "", logData, log.Error(ErrorInvalidHealthStatus))
	}
	return cli.serviceName, ErrorInvalidHealthStatus
}

// Checker : Check health of Elasticsearch and return it inside a Check structure
func (cli *Client) Checker(ctx *context.Context) (*health.Check, error) {
	_, err := cli.healthcheck()
	if err != nil {
		switch err {
		case ErrorClusterAtRisk:
			return cli.getCheck(ctx, 429), err
		default:
			return cli.getCheck(ctx, 500), err
		}
	}
	return cli.getCheck(ctx, 200), nil
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

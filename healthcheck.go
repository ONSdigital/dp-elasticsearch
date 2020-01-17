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

// MsgHealthy Check message returned when elasticsearch is healthy
const MsgHealthy = "elasticsearch is healthy"

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

// minTime is the oldest time for Check structure.
var minTime = time.Unix(0, 0)

// ClusterHealth represents the response from the elasticsearch cluster health check
type ClusterHealth struct {
	Status string `json:"status"`
}

// healthcheck calls elasticsearch to check its health status. This call implements only the logic,
// without providing the Check object, and it's aimed for internal use.
func (cli *Client) healthcheck(ctx context.Context) (code int, err error) {

	urlHealth := cli.url + pathHealth
	logData := log.Data{"url": urlHealth}

	URL, err := url.Parse(urlHealth)
	if err != nil {
		log.Event(nil, "failed to create url for elasticsearch healthcheck", logData, log.Error(err))
		return 500, err
	}

	path := URL.String()
	logData["url"] = path

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Event(nil, "failed to create request for healthcheck call to elasticsearch", logData, log.Error(err))
		return 500, err
	}

	if cli.signRequests {
		awsauth.Sign(req)
	}

	resp, err := cli.httpCli.Do(ctx, req)
	if err != nil {
		log.Event(nil, "failed to call elasticsearch", logData, log.Error(err))
		return 500, err
	}
	defer resp.Body.Close()

	logData["http_code"] = resp.StatusCode
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		log.Event(nil, "", logData, log.Error(ErrorUnexpectedStatusCode))
		return resp.StatusCode, ErrorUnexpectedStatusCode
	}

	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Event(nil, "failed to read response body from call to elastic", logData, log.Error(err))
		return resp.StatusCode, ErrorUnexpectedStatusCode
	}

	var clusterHealth ClusterHealth
	err = json.Unmarshal(jsonBody, &clusterHealth)
	if err != nil {
		log.Event(nil, "", logData, log.Error(ErrorParsingBody))
		return resp.StatusCode, ErrorParsingBody
	}

	logData["cluster_health"] = clusterHealth.Status
	switch clusterHealth.Status {
	case healthValues[HealthGreen]:
		return resp.StatusCode, nil
	case healthValues[HealthYellow]:
		log.Event(nil, "", logData, log.Error(ErrorClusterAtRisk))
		return resp.StatusCode, ErrorClusterAtRisk
	case healthValues[HealthRed]:
		log.Event(nil, "", logData, log.Error(ErrorUnhealthyClusterStatus))
		return resp.StatusCode, ErrorUnhealthyClusterStatus
	default:
		log.Event(nil, "", logData, log.Error(ErrorInvalidHealthStatus))
	}
	return resp.StatusCode, ErrorInvalidHealthStatus
}

// Checker checks health of Elasticsearch and return it inside a Check structure. This method decides the severity of any possible error.
func (cli *Client) Checker(ctx *context.Context) (*health.Check, error) {
	statusCode, err := cli.healthcheck(*ctx)
	currentTime := time.Now().UTC()
	cli.Check.LastChecked = &currentTime
	cli.Check.StatusCode = statusCode
	if err != nil {
		cli.Check.LastFailure = &currentTime
		cli.Check.Status = getStatusFromError(err)
		cli.Check.Message = err.Error()
		return cli.Check, err
	}
	cli.Check.LastSuccess = &currentTime
	cli.Check.Status = health.StatusOK
	cli.Check.Message = MsgHealthy
	return cli.Check, nil
}

// getStatusFromError decides the health status (severity) according to the provided error
func getStatusFromError(err error) string {
	switch err {
	case ErrorClusterAtRisk:
		return health.StatusWarning
	default:
		return health.StatusCritical
	}
}

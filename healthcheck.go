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

// MsgHealthy Check message returned when elasticsearch is healthy and indexed
const MsgHealthy = "elasticsearch is indexed and healthy"

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
	ErrorUnexpectedStatusCode    = errors.New("unexpected status code from api")
	ErrorParsingBody             = errors.New("error parsing cluster health response body")
	ErrorClusterAtRisk           = errors.New("error cluster state yellow. Data might be at risk, check your replica shards")
	ErrorUnhealthyClusterStatus  = errors.New("error cluster health red. Cluster is unhealthy")
	ErrorInvalidHealthStatus     = errors.New("error invalid health status returned")
	ErrorIndexDoesNotExist       = errors.New("error index does not exist in cluster")
	ErrorRetrievingClientIndices = errors.New("error retrieving indices from client")
)

// minTime is the oldest time for Check structure.
var minTime = time.Unix(0, 0)

// ClusterHealth represents the response from the elasticsearch cluster health check
type ClusterHealth struct {
	Status string `json:"status"`
}

//indexcheck calls elasticsearch to check if an index from the client exists
//What is a context?
func (cli *Client) indexcheck(ctx context.Context) (code int, err error) {

	var clientElasticSearchIndices = cli.indices

	for _, index := range clientElasticSearchIndices {

		urlIndex := cli.url + "/" + index
		logData := log.Data{"url": urlIndex}

		URL, err := url.Parse(urlIndex)
		if err != nil {
			log.Event(ctx, "failed to create url for elasticsearch indexcheck", log.ERROR, logData, log.Error(err))
			return 500, err
		}

		path := URL.String()
		//What do the ["url"] brackets mean - is it resetting the value for the "url" key to the path
		logData["url"] = path

		//what is background context?
		req, err := http.NewRequest("HEAD", path, nil)
		if err != nil {
			log.Event(ctx, "failed to create request for indexcheck call to elasticsearch", log.ERROR, logData, log.Error(err))
			return 500, err
		}

		//Where do the credentials come from?
		if cli.signRequests {
			awsauth.Sign(req)
		}

		resp, err := cli.httpCli.Do(ctx, req)
		if err != nil {
			log.Event(ctx, "failed to call elasticsearch", log.ERROR, logData, log.Error(err))
			return 500, err
		}
		defer resp.Body.Close()
		logData["http_code"] = resp.StatusCode

		switch resp.StatusCode {

		case 200:
			return resp.StatusCode, nil
		case 404:
			log.Event(ctx, "404 status code returned in response", logData, log.ERROR, log.Error(ErrorIndexDoesNotExist))
			return resp.StatusCode, ErrorIndexDoesNotExist
		default:
			log.Event(ctx, "unexpected status code returned in response", logData, log.ERROR, log.Error(ErrorUnexpectedStatusCode))
			return resp.StatusCode, ErrorUnexpectedStatusCode
		}

	}

	//I think I need to return something here (there is an error if not) - just put in this to see if it got rid of the error
	//Would it return here if it was unable to form the var clientElasticSearchIndices = cli.indices because it couldn't find the cli.indices?
	//I also couldn't work out how to wrap the return and the log in an if err statement because if I did there was still a need to return something

	logData := log.Data{"clientIndexResponse": cli.indices}
	log.Event(ctx, "failed to retrieve the indices from the client", log.ERROR, logData, log.Error(ErrorRetrievingClientIndices))
	return 500, ErrorRetrievingClientIndices

}

// healthcheck calls elasticsearch to check its health status. This call implements only the logic,
// without providing the Check object, and it's aimed for internal use.
func (cli *Client) healthcheck(ctx context.Context) (code int, err error) {

	urlHealth := cli.url + pathHealth
	logData := log.Data{"url": urlHealth}

	URL, err := url.Parse(urlHealth)
	if err != nil {
		log.Event(ctx, "failed to create url for elasticsearch healthcheck", log.ERROR, logData, log.Error(err))
		return 500, err
	}

	path := URL.String()
	logData["url"] = path

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Event(ctx, "failed to create request for healthcheck call to elasticsearch", log.ERROR, logData, log.Error(err))
		return 500, err
	}

	if cli.signRequests {
		awsauth.Sign(req)
	}

	resp, err := cli.httpCli.Do(ctx, req)
	if err != nil {
		log.Event(ctx, "failed to call elasticsearch", log.ERROR, logData, log.Error(err))
		return 500, err
	}
	defer resp.Body.Close()

	logData["http_code"] = resp.StatusCode
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		log.Event(ctx, "unexpected status code returned in response", logData, log.ERROR, log.Error(ErrorUnexpectedStatusCode))
		return resp.StatusCode, ErrorUnexpectedStatusCode
	}

	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Event(ctx, "failed to read response body from call to elastic", log.ERROR, logData, log.Error(err))
		return resp.StatusCode, ErrorUnexpectedStatusCode
	}

	var clusterHealth ClusterHealth
	err = json.Unmarshal(jsonBody, &clusterHealth)
	if err != nil {
		log.Event(ctx, "json unmarshal error", log.ERROR, logData, log.Error(ErrorParsingBody))
		return resp.StatusCode, ErrorParsingBody
	}

	logData["cluster_health"] = clusterHealth.Status
	switch clusterHealth.Status {
	case healthValues[HealthGreen]:
		return resp.StatusCode, nil
	case healthValues[HealthYellow]:
		log.Event(ctx, "yellow health status", log.WARN, logData, log.Error(ErrorClusterAtRisk))
		return resp.StatusCode, ErrorClusterAtRisk
	case healthValues[HealthRed]:
		log.Event(ctx, "red health status", log.WARN, logData, log.Error(ErrorUnhealthyClusterStatus))
		return resp.StatusCode, ErrorUnhealthyClusterStatus
	default:
		log.Event(ctx, "invalid health status", log.WARN, logData, log.Error(ErrorInvalidHealthStatus))
	}
	return resp.StatusCode, ErrorInvalidHealthStatus
}

// Checker checks health of Elasticsearch and if there is an index and updates the provided CheckState accordingly.
func (cli *Client) Checker(ctx context.Context, state *health.CheckState) error {
	statusCodeHealthCheck, err := cli.healthcheck(ctx)
	if err != nil {
		state.Update(getStatusFromError(err), err.Error(), statusCodeHealthCheck)
		return nil
	}

	statusCodeIndexCheck, err := cli.indexcheck(ctx)
	if err != nil {
		state.Update(getStatusFromError(err), err.Error(), statusCodeIndexCheck)
		return nil
	}

	state.Update(health.StatusOK, MsgHealthy, statusCodeIndexCheck)
	return nil
}

// getStatusFromError decides the health status (severity) according to the provided error.
func getStatusFromError(err error) string {
	switch err {
	case ErrorClusterAtRisk:
		return health.StatusWarning
	default:
		return health.StatusCritical
	}
}

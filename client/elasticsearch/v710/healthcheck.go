package v710

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/v2/log"
)

// MsgHealthy Check message returned when elasticsearch is healthy and the required indexes exist
const MsgHealthy = "elasticsearch is healthy and the required indexes exist"

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
	ErrorClusterAtRisk          = errors.New("elasticsearch cluster state yellow but functional. Data might be at risk, check your replica shards")
	ErrorUnhealthyClusterStatus = errors.New("error cluster health red. Cluster is unhealthy")
	ErrorInvalidHealthStatus    = errors.New("error invalid health status returned")
	ErrorIndexDoesNotExist      = errors.New("error index does not exist in cluster")
	ErrorInternalServer         = errors.New("error internal server error")
)

// ClusterHealth represents the response from the elasticsearch cluster health check
type ClusterHealth struct {
	Status string `json:"status"`
}

// Checker checks health of Elasticsearch, if the required indexes exist and updates the provided CheckState accordingly.
func (cli *ESClient) Checker(ctx context.Context, state *health.CheckState) error {
	if state == nil {
		state = &health.CheckState{}
	}

	statusCode, err := cli.healthcheck(ctx)
	if err != nil && err != ErrorClusterAtRisk {
		if updateErr := state.Update(health.StatusCritical, err.Error(), statusCode); err != nil {
			log.Warn(ctx, "unable to update health state", log.FormatErrors([]error{updateErr}))
		}

		return nil
	}

	if len(cli.indexes) > 0 {
		if indexStatusCode, indexErr := cli.indexcheck(ctx); indexErr != nil {
			if updateErr := state.Update(health.StatusCritical, indexErr.Error(), indexStatusCode); err != nil {
				log.Warn(ctx, "unable to update health state", log.FormatErrors([]error{updateErr}))
			}

			return nil
		}
	}

	// Elasticsearch cluster configuration should not determine if the health check should fail
	// The application will still be able to communicate to the elasticsearch cluster - hence es
	// responding with 200 staus code in response
	if err == ErrorClusterAtRisk {
		if updateErr := state.Update(health.StatusOK, err.Error(), statusCode); err != nil {
			log.Warn(ctx, "unable to update health state", log.FormatErrors([]error{updateErr}))
		}

		return nil
	}

	if updateErr := state.Update(health.StatusOK, MsgHealthy, statusCode); err != nil {
		log.Warn(ctx, "unable to update health state", log.FormatErrors([]error{updateErr}))
	}

	return nil
}

// healthcheck calls elasticsearch to check its health status. This call implements only the logic,
// without providing the Check object, and it's aimed for internal use.
func (cli *ESClient) healthcheck(ctx context.Context) (code int, err error) {
	resp, err := cli.esClient.Cluster.Health()
	if err != nil {
		log.Error(ctx, "failed to call elasticsearch", err)
		return 500, err
	}
	defer resp.Body.Close()

	logData := log.Data{"http_code": resp.StatusCode}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		log.Error(ctx, "unexpected status code returned in response", ErrorUnexpectedStatusCode)
		return resp.StatusCode, ErrorUnexpectedStatusCode
	}

	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(ctx, "failed to read response body from call to elastic", err)
		return resp.StatusCode, ErrorUnexpectedStatusCode
	}

	var clusterHealth ClusterHealth
	err = json.Unmarshal(jsonBody, &clusterHealth)
	if err != nil {
		log.Error(ctx, "json unmarshal error", ErrorParsingBody)
		return resp.StatusCode, ErrorParsingBody
	}

	logData["cluster_health"] = clusterHealth.Status
	switch clusterHealth.Status {
	case healthValues[HealthGreen]:
		return resp.StatusCode, nil
	case healthValues[HealthYellow]:
		log.Error(ctx, "yellow health status", ErrorClusterAtRisk)
		return resp.StatusCode, ErrorClusterAtRisk
	case healthValues[HealthRed]:
		log.Error(ctx, "red health status", ErrorUnhealthyClusterStatus)
		return resp.StatusCode, ErrorUnhealthyClusterStatus
	default:
		log.Error(ctx, "invalid health status", ErrorInvalidHealthStatus)
	}

	return resp.StatusCode, ErrorInvalidHealthStatus
}

// indexcheck calls elasticsearch to check if the required indexes from the client exist
func (cli *ESClient) indexcheck(ctx context.Context) (code int, err error) {
	for _, index := range cli.indexes {
		resp, err := cli.esClient.Cluster.Health(cli.esClient.Cluster.Health.WithIndex(index))
		if err != nil {
			log.Error(ctx, "failed to call elasticsearch", err)
			return 500, err
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case 200:
			continue
		case 404:
			log.Error(ctx, "index does not exist", ErrorIndexDoesNotExist)
			return resp.StatusCode, ErrorIndexDoesNotExist
		default:
			log.Error(ctx, "unexpected status code returned in response", ErrorUnexpectedStatusCode)
			return resp.StatusCode, ErrorUnexpectedStatusCode
		}
	}

	return 200, nil
}

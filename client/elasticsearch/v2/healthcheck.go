package v2

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/v2/log"
)

// HTTP path to check the health of a cluster using Elasticsearch API
const pathHealth = "/_cluster/health"

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

// indexcheck calls elasticsearch to check if the required indexes from the client exist
func (cli *Client) indexcheck(ctx context.Context) (code int, err error) {
	// Check handles each index, making sure the response body is always closed
	check := func(index string) (int, error) {
		urlIndex := cli.url + "/" + index

		_, err := url.Parse(urlIndex)
		if err != nil {
			log.Error(ctx, "failed to create url for elasticsearch indexcheck", err)
			return 500, err
		}

		req, err := http.NewRequest("HEAD", urlIndex, http.NoBody)
		if err != nil {
			log.Error(ctx, "failed to create request for indexcheck call to elasticsearch", err)
			return 500, err
		}

		resp, err := cli.httpCli.Do(ctx, req)
		if err != nil {
			log.Error(ctx, "failed to call elasticsearch", err)
			return 500, err
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case 200:
			return 200, nil
		case 404:
			log.Error(ctx, "index does not exist", ErrorIndexDoesNotExist)
			return resp.StatusCode, ErrorIndexDoesNotExist
		default:
			log.Error(ctx, "unexpected status code returned in response", ErrorUnexpectedStatusCode)
			return resp.StatusCode, ErrorUnexpectedStatusCode
		}
	}

	// Check all indexes, if any fails, return the code and error
	for _, index := range cli.indexes {
		code, err := check(index)
		if err != nil {
			return code, err
		}
	}

	// if all indexes are successful, return 200 and no error
	return 200, nil
}

// healthcheck calls elasticsearch to check its health status. This call implements only the logic,
// without providing the Check object, and it's aimed for internal use.
func (cli *Client) healthcheck(ctx context.Context) (code int, err error) {
	urlHealth := cli.url + pathHealth
	logData := log.Data{"url": urlHealth}

	URL, err := url.Parse(urlHealth)
	if err != nil {
		log.Error(ctx, "failed to create url for elasticsearch healthcheck", err)
		return 500, err
	}

	path := URL.String()
	logData["url"] = path

	req, err := http.NewRequest("GET", path, http.NoBody)
	if err != nil {
		log.Error(ctx, "failed to create request for healthcheck call to elasticsearch", err)
		return 500, err
	}

	resp, err := cli.httpCli.Do(ctx, req)
	if err != nil {
		log.Error(ctx, "failed to call elasticsearch", err)
		return 500, err
	}
	defer resp.Body.Close()

	logData["http_code"] = resp.StatusCode
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		log.Error(ctx, "unexpected status code returned in response", ErrorUnexpectedStatusCode)
		return resp.StatusCode, ErrorUnexpectedStatusCode
	}

	jsonBody, err := io.ReadAll(resp.Body)
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

// Checker checks health of Elasticsearch, if the required indexes exist and updates the provided CheckState accordingly.
func (cli *Client) Checker(ctx context.Context, state *health.CheckState) error {
	if state == nil {
		state = &health.CheckState{}
	}

	statusCode, err := cli.healthcheck(ctx)
	if err != nil && err != ErrorClusterAtRisk {
		if updateErr := state.Update(health.StatusCritical, err.Error(), statusCode); updateErr != nil {
			log.Warn(ctx, "unable to update health state", log.FormatErrors([]error{updateErr}))
		}

		return nil
	}

	if len(cli.indexes) > 0 {
		if indexStatusCode, indexErr := cli.indexcheck(ctx); indexErr != nil {
			if updateErr := state.Update(health.StatusCritical, indexErr.Error(), indexStatusCode); updateErr != nil {
				log.Warn(ctx, "unable to update health state", log.FormatErrors([]error{updateErr}))
			}

			return nil
		}
	}

	// Elasticsearch cluster configuration should not determine if the health check should fail
	// The application will still be able to communicate to the elasticsearch cluster - hence es
	// responding with 200 staus code in response
	if err == ErrorClusterAtRisk {
		if updateErr := state.Update(health.StatusOK, err.Error(), statusCode); updateErr != nil {
			log.Warn(ctx, "unable to update health state", log.FormatErrors([]error{updateErr}))
		}

		return nil
	}

	if updateErr := state.Update(health.StatusOK, MsgHealthy, statusCode); updateErr != nil {
		log.Warn(ctx, "unable to update health state", log.FormatErrors([]error{updateErr}))
	}

	return nil
}

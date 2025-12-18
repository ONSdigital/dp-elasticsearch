package elasticsearch

import (
	"fmt"

	"github.com/ONSdigital/dp-elasticsearch/v3/client"
	v710 "github.com/ONSdigital/dp-elasticsearch/v3/client/elasticsearch/v710"
)

func NewClient(cfg client.Config) (client.Client, error) {
	switch cfg.ClientLib {
	case client.OpenSearch:
		return nil, fmt.Errorf("the Opensearch client is currently not implemented")
	default:
		return v710.NewESClient(cfg.Address, cfg.Transport)
	}
}

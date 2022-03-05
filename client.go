package elasticsearch

import (
	"fmt"

	"github.com/ONSdigital/dp-elasticsearch/v3/client"
	v2 "github.com/ONSdigital/dp-elasticsearch/v3/client/elasticsearch/v2"
	v710 "github.com/ONSdigital/dp-elasticsearch/v3/client/elasticsearch/v710"
)

func NewClient(cfg client.Config) (client.Client, error) {
	switch cfg.ClientLib {
	case client.GoElasticV710:
		return v710.NewESClient(cfg.Address, cfg.Transport)
	case client.OpenSearch:
		return nil, fmt.Errorf("the Opensearch client is currently not implemented")
	default:
		return v2.NewClient(cfg.Address, cfg.MaxRetries, cfg.Indexes...), nil
	}
}

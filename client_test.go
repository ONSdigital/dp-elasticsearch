package elasticsearch_test

import (
	"github.com/ONSdigital/dp-elasticsearch/v3"
	"github.com/ONSdigital/dp-elasticsearch/v3/clients/elasticsearch/v2"
	"github.com/ONSdigital/dp-elasticsearch/v3/clients/elasticsearch/v710"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewClient_ReturnsNewGoElasticClientVersion710(t *testing.T) {
	t.Parallel()
	cfg := elasticsearch.Config{
		ClientLib: elasticsearch.GoElastic_V710,
		Address:   "http://some-url.com",
	}

	cli, err := elasticsearch.NewClient(cfg)

	assert.Nil(t, err)
	assert.NotNil(t, cli)
	assert.IsType(t, cli, &v710.ESClient{})
}

func TestNewClient_WhenValidURLIsNotSpecified_ReturnsError(t *testing.T) {
	t.Parallel()
	cfg := elasticsearch.Config{
		ClientLib: elasticsearch.GoElastic_V710,
		Address:   "invalid-url",
	}

	cli, err := elasticsearch.NewClient(cfg)

	assert.NotNil(t, err)
	assert.Nil(t, cli)
}

func TestNewClient_WhenOpenSearchClientLibraryIsRequested_ReturnsNotImplemented(t *testing.T) {
	t.Parallel()
	cfg := elasticsearch.Config{
		ClientLib: elasticsearch.OpenSearch,
		Address:   "http://some-url.com",
	}

	cli, err := elasticsearch.NewClient(cfg)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "not implemented")
	assert.Nil(t, cli)
}

func TestNewClient_ReturnsNewDefaultClient(t *testing.T) {
	t.Parallel()
	cfg := elasticsearch.Config{
		Address: "http://some-url.com",
	}

	cli, err := elasticsearch.NewClient(cfg)

	assert.Nil(t, err)
	assert.NotNil(t, cli)
	assert.IsType(t, cli, &v2.Client{})
}

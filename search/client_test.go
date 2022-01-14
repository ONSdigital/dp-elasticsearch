package search_test

import (
	"testing"

	"github.com/ONSdigital/dp-elasticsearch/v2/elasticsearch"
	"github.com/ONSdigital/dp-elasticsearch/v2/go_elasticsearch_v710"
	"github.com/ONSdigital/dp-elasticsearch/v2/search"
	"github.com/stretchr/testify/assert"
)

func TestNewClient_ReturnsNewGoElasticClientVersion710(t *testing.T) {
	t.Parallel()
	cfg := search.Config{
		ClientLib: search.GoElastic_V710,
		Address:   "http://some-url.com",
		IndexName: "some index",
	}

	cli, err := search.NewClient(cfg)

	assert.Nil(t, err)
	assert.NotNil(t, cli)
	assert.IsType(t, cli, &go_elasticsearch_v710.ESClient{})
}

func TestNewClient_WhenIndexNameIsNotSpecified_ReturnsError(t *testing.T) {
	t.Parallel()
	cfg := search.Config{
		ClientLib: search.GoElastic_V710,
		Address:   "http://some-url.com",
	}

	cli, err := search.NewClient(cfg)

	assert.NotNil(t, err)
	assert.Nil(t, cli)
}

func TestNewClient_WhenValidURLIsNotSpecified_ReturnsError(t *testing.T) {
	t.Parallel()
	cfg := search.Config{
		ClientLib: search.GoElastic_V710,
		Address:   "invalid-url",
	}

	cli, err := search.NewClient(cfg)

	assert.NotNil(t, err)
	assert.Nil(t, cli)
}

func TestNewClient_WhenOpenSearchClientLibraryIsRequested_ReturnsNotImplemented(t *testing.T) {
	t.Parallel()
	cfg := search.Config{
		ClientLib: search.OpenSearch,
		Address:   "http://some-url.com",
	}

	cli, err := search.NewClient(cfg)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "not implemented")
	assert.Nil(t, cli)
}

func TestNewClient_ReturnsNewDefaultClient(t *testing.T) {
	t.Parallel()
	cfg := search.Config{
		Address: "http://some-url.com",
	}

	cli, err := search.NewClient(cfg)

	assert.Nil(t, err)
	assert.NotNil(t, cli)
	assert.IsType(t, cli, &elasticsearch.Client{})
}

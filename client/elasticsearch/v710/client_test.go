package v710

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-elasticsearch/v4/client"
	es710 "github.com/elastic/go-elasticsearch/v7"
	. "github.com/smartystreets/goconvey/convey"
)

type mockRoundTripper struct {
	roundTripFunc func(req *http.Request) *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req), nil
}

func newMockClient(statusCode int, body string, assert func(req *http.Request)) *es710.Client {
	rt := &mockRoundTripper{
		roundTripFunc: func(req *http.Request) *http.Response {
			if assert != nil {
				assert(req)
			}
			return &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(bytes.NewBufferString(body)),
				Header:     make(http.Header),
			}
		},
	}
	es, _ := es710.NewClient(es710.Config{
		Addresses: []string{"http://localhost:9200"},
		Transport: rt,
	})
	return es
}

func TestMultiSearch(t *testing.T) {
	t.Parallel()

	Convey("Given convert a slice of searches to multiline searches", t, func() {
		expectedMultiLintStringCount := 5
		searches := []client.Search{
			{
				Header: client.Header{
					Index: "ons_test",
				},
				Query: []byte(`{"query" : {"match" : { "message": "this is a test"}}}`),
			},
			{
				Header: client.Header{
					Index: "ons_test_2",
				},
				Query: []byte(`{"query" : {"match_all" : {}}}`),
			},
		}

		body, err := convertToMultilineSearches(searches)

		So(err, ShouldEqual, nil)
		splitQuery := strings.Split(string(body), "\n")
		So(len(splitQuery), ShouldEqual, expectedMultiLintStringCount)
		So(splitQuery[0], ShouldEqual, "{\"index\":\"ons_test\"}")
		So(splitQuery[1], ShouldEqual, "{\"query\" : {\"match\" : { \"message\": \"this is a test\"}}}")
		So(splitQuery[2], ShouldEqual, "{\"index\":\"ons_test_2\"}")
		So(splitQuery[3], ShouldEqual, "{\"query\" : {\"match_all\" : {}}}")
	})
}

func TestDeleteDocument(t *testing.T) {
	Convey("Given a valid ESClient", t, func() {
		esClient := newMockClient(http.StatusOK, `{}`, nil)
		testClient := &ESClient{esClient: esClient}

		Convey("When DeleteDocument returns 200", func() {
			err := testClient.DeleteDocument(context.Background(), "my-index", "my-id")
			So(err, ShouldBeNil)
		})

		Convey("When DeleteDocument returns 500", func() {
			esClient := newMockClient(http.StatusInternalServerError, `{"error":"server error"}`, nil)
			testClient := &ESClient{esClient: esClient}
			err := testClient.DeleteDocument(context.Background(), "my-index", "my-id")
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "delete request failed")
		})
	})
}

func TestDeleteDocumentByQuery(t *testing.T) {
	Convey("Given a valid ESClient", t, func() {
		var receivedBody string
		assertBody := func(req *http.Request) {
			bodyBytes, _ := io.ReadAll(req.Body)
			receivedBody = string(bodyBytes)
		}

		esClient := newMockClient(http.StatusOK, `{}`, assertBody)
		testClient := &ESClient{esClient: esClient}

		Convey("When DeleteDocumentByQuery sends a correct query", func() {
			query := `{
				"query": {
					"term": {
						"uri": "/my/uri"
					}
				}
			}`

			search := client.Search{
				Header: client.Header{Index: "my-index"},
				Query:  []byte(query),
			}

			err := testClient.DeleteDocumentByQuery(context.Background(), search)
			So(err, ShouldBeNil)
			So(receivedBody, ShouldContainSubstring, `"term"`)
			So(receivedBody, ShouldContainSubstring, `"/my/uri"`)
		})

		Convey("When DeleteDocumentByQuery returns 500", func() {
			errorClient := newMockClient(http.StatusInternalServerError, `{"error":"bad stuff"}`, nil)
			testClient := &ESClient{esClient: errorClient}

			query := `{
				"query": {
					"term": {
						"uri": "/bad/uri"
					}
				}
			}`

			search := client.Search{
				Header: client.Header{Index: "my-index"},
				Query:  []byte(query),
			}

			err := testClient.DeleteDocumentByQuery(context.Background(), search)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "delete-by-query failed")
		})
	})
}

package v710

import (
	"strings"
	"testing"

	"github.com/ONSdigital/dp-elasticsearch/v3/client"

	. "github.com/smartystreets/goconvey/convey"
)

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

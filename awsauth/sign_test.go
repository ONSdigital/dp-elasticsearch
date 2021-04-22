package awsauth

import (
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateNewSigner(t *testing.T) {
	Convey("Given that we want to use the aws sdk signer", t, func() {
		Convey("When the region is set to an empty string", func() {
			signer := NewSigner(true, "", "es")

			Convey("Then an error is returned when attempting to Sign the request", func() {
				req := httptest.NewRequest("GET", "http://test-url", nil)

				err := signer.Sign(req, nil, time.Now())
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, errors.New("No AWS region was provided. Cannot sign request."))
			})
		})

		Convey("When the service is set to an empty string", func() {
			signer := NewSigner(true, "eu-west-1", "")

			Convey("Then an error is returned when attempting to Sign the request", func() {
				req := httptest.NewRequest("GET", "http://test-url", nil)

				err := signer.Sign(req, nil, time.Now())
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, errors.New("No AWS service was provided. Cannot sign request."))
			})
		})

		// Convey("When the service and region are set", func() {
		// 	signer := NewSigner(true, "eu-west-1", "es")

		// 	Convey("Then no error is returned when attempting to Sign the request", func() {
		// 		req := httptest.NewRequest("GET", "http://test-url", nil)

		// 		err := signer.Sign(req, nil, time.Now())
		// 		So(err, ShouldBeNil)
		// 	})
		// })
	})

	Convey("Given that we want to use the smartystreets auth signage package", t, func() {
		Convey("When the service and region are not set", func() {
			signer := NewSigner(false, "", "")

			Convey("Then no error is returned when attempting to Sign the request", func() {
				req := httptest.NewRequest("GET", "http://test-url", nil)

				err := signer.Sign(req, nil, time.Now())
				So(err, ShouldBeNil)
			})
		})
	})
}

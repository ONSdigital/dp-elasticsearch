package awsauth

import (
	"errors"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	testAccessKey       = "TEST_ACCESS_KEY"
	testSecretAccessKey = "TEST_SECRET_KEY"
)

func TestCreateNewSigner(t *testing.T) {
	removeTestEnvironmentVariables()

	Convey("Given that we want to use the aws sdk signer", t, func() {
		Convey("When the region is set to an empty string", func() {
			signer := NewSigner(true, "", "", "es")

			Convey("Then an error is returned when attempting to Sign the request", func() {
				req := httptest.NewRequest("GET", "http://test-url", nil)

				err := signer.Sign(req, nil, time.Now())
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, errors.New("No AWS region was provided. Cannot sign request."))
			})
		})

		Convey("When the service is set to an empty string", func() {
			signer := NewSigner(true, "", "eu-west-1", "")

			Convey("Then an error is returned when attempting to Sign the request", func() {
				req := httptest.NewRequest("GET", "http://test-url", nil)

				err := signer.Sign(req, nil, time.Now())
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, errors.New("No AWS service was provided. Cannot sign request."))
			})
		})

		Convey("When the service and region are set and credentials are set in environment variables", func() {
			os.Setenv(envAccessKeyID, testAccessKey)
			os.Setenv(envSecretAccessKey, testSecretAccessKey)

			signer := NewSigner(true, "", "eu-west-1", "es")

			Convey("Then no error is returned when attempting to Sign the request", func() {
				req := httptest.NewRequest("GET", "http://test-url", nil)

				err := signer.Sign(req, nil, time.Now())
				So(err, ShouldBeNil)
			})

			removeTestEnvironmentVariables()
		})
	})

	Convey("Given that we want to use the smartystreets auth signage package", t, func() {
		Convey("When the service and region are not set", func() {
			signer := NewSigner(false, "", "", "")

			Convey("Then no error is returned when attempting to Sign the request", func() {
				req := httptest.NewRequest("GET", "http://test-url", nil)

				err := signer.Sign(req, nil, time.Now())
				So(err, ShouldBeNil)
			})
		})
	})
}

func removeTestEnvironmentVariables() {
	os.Setenv(envAccessKey, "")
	os.Setenv(envAccessKeyID, "")
	os.Setenv(envSecretKey, "")
	os.Setenv(envSecretAccessKey, "")
	os.Setenv(envSessionToken, "")
}

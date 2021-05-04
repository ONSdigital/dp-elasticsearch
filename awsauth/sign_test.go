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
	envAccessKey       = "AWS_ACCESS_KEY"
	envAccessKeyID     = "AWS_ACCESS_KEY_ID"
	envSecretKey       = "AWS_SECRET_KEY"
	envSecretAccessKey = "AWS_SECRET_ACCESS_KEY"

	testAccessKey       = "TEST_ACCESS_KEY"
	testSecretAccessKey = "TEST_SECRET_KEY"
)

func TestCreateNewSigner(t *testing.T) {
	removeTestEnvironmentVariables()

	Convey("Given that we want to use the aws sdk signer", t, func() {
		Convey("When the region is set to an empty string", func() {
			Convey("Then an error is returned when retrieving aws sdk signer", func() {
				signer, err := NewAwsSigner("", "", "", "es")
				So(err, ShouldResemble, errors.New("No AWS region was provided. Cannot sign request."))
				So(signer, ShouldBeNil)

				Convey("But no error is returned when attempting to Sign the request due to fallback to using smartystreets signer", func() {
					req := httptest.NewRequest("GET", "http://test-url", nil)

					err = signer.Sign(req, nil, time.Now())
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("When the service is set to an empty string", func() {
			Convey("Then an error is returned when retrieving aws sdk signer", func() {
				signer, err := NewAwsSigner("", "", "eu-west-1", "")
				So(err, ShouldResemble, errors.New("No AWS service was provided. Cannot sign request."))
				So(signer, ShouldBeNil)

				Convey("But no error is returned when attempting to Sign the request due to fallback to using smartystreets signer", func() {
					req := httptest.NewRequest("GET", "http://test-url", nil)

					err = signer.Sign(req, nil, time.Now())
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("When the service and region are set and credentials are set in environment variables", func() {
			os.Setenv(envAccessKeyID, testAccessKey)
			os.Setenv(envSecretAccessKey, testSecretAccessKey)

			Convey("Then an error is returned when retrieving aws sdk signer", func() {
				signer, err := NewAwsSigner("", "", "eu-west-1", "es")
				So(err, ShouldBeNil)
				So(signer, ShouldNotBeNil)

				Convey("And no error is returned when attempting to Sign the request", func() {
					req := httptest.NewRequest("GET", "http://test-url", nil)

					err := signer.Sign(req, nil, time.Now())
					So(err, ShouldBeNil)
				})
			})

			removeTestEnvironmentVariables()
		})
	})

	Convey("Given that we want to use the smartystreets auth signage package", t, func() {
		Convey("When the service and region are not set", func() {
			signer := &Signer{}

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
}

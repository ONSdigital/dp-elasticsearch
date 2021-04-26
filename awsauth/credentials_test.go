package awsauth

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	testPrimaryAccessKey         = "TEST_PRIMARY_ACCESS_KEY"
	testPrimarySecretAccessKey   = "TEST_PRIMARY_SECRET_KEY"
	testSecondaryAccessKey       = "TEST_SECONDARY_ACCESS_KEY"
	testSecondarySecretAccessKey = "TEST_SECONDARY_SECRET_KEY"
	testToken                    = "TESTTOKEN"
)

func TestRetrievalOfCredentials(t *testing.T) {
	// Clear down local environment variables before running test
	removeTestEnvironmentVariables()

	Convey("Given a credentials provider", t, func() {
		provider := getProvider()

		Convey("When attempting to retrieve credentials with primary environment variable set", func() {
			os.Setenv(envAccessKeyID, testPrimaryAccessKey)
			os.Setenv(envSecretAccessKey, testPrimarySecretAccessKey)
			os.Setenv(envSessionToken, testToken)

			creds, err := provider.Retrieve()

			Convey("Then credential values are returned successfully", func() {
				So(err, ShouldBeNil)
				So(creds.AccessKeyID, ShouldEqual, testPrimaryAccessKey)
				So(creds.SecretAccessKey, ShouldEqual, testPrimarySecretAccessKey)
				So(creds.SessionToken, ShouldEqual, testToken)
			})

			Convey("And that the credentials are not expired", func() {
				So(provider.IsExpired(), ShouldBeFalse)
			})

			removeTestEnvironmentVariables()
		})

		Convey("When attempting to retrieve credentials with secondary set of environment variables", func() {
			os.Setenv(envAccessKey, testSecondaryAccessKey)
			os.Setenv(envSecretKey, testSecondarySecretAccessKey)
			os.Setenv(envSessionToken, testToken)

			creds, err := provider.Retrieve()

			Convey("Then credential values are returned successfully", func() {
				So(err, ShouldBeNil)
				So(creds.AccessKeyID, ShouldEqual, testSecondaryAccessKey)
				So(creds.SecretAccessKey, ShouldEqual, testSecondarySecretAccessKey)
				So(creds.SessionToken, ShouldEqual, testToken)
			})

			Convey("And that the credentials are not expired", func() {
				So(provider.IsExpired(), ShouldBeFalse)
			})

			removeTestEnvironmentVariables()
		})
	})
}

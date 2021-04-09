package awsauth

import (
	"errors"
	"io"
	"net/http"
	"time"

	signerV4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	awsauth "github.com/smartystreets/go-aws-auth"
)

type Signer struct {
	awsRegion    string
	awsSDKSigner bool
	awsService   string
}

func NewSigner(awsSDKSigner bool, awsRegion, awsService string) *Signer {
	return &Signer{
		awsRegion:    awsRegion,
		awsSDKSigner: awsSDKSigner,
		awsService:   awsService,
	}
}

func (signer *Signer) Sign(req *http.Request, bodyReader io.ReadSeeker, currentTime time.Time) error {
	if signer.awsSDKSigner {
		if err := signer.validateAwsSDKSigner(); err != nil {
			return err
		}

		creds := retrieveCredentials()
		v4Signer := signerV4.NewSigner(creds)
		v4Signer.Sign(req, bodyReader, signer.awsService, signer.awsRegion, time.Now())
	} else {
		awsauth.Sign(req)
	}

	return nil
}

func (signer *Signer) validateAwsSDKSigner() error {
	if signer.awsRegion == "" {
		return errors.New("No AWS region was provided. Cannot sign request.")
	}

	if signer.awsService == "" {
		return errors.New("No AWS service was provided. Cannot sign request.")
	}

	return nil
}

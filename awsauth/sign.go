package awsauth

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	signerV4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	awsauth "github.com/smartystreets/go-aws-auth"
)

type Signer struct {
	awsRegion  string
	awsService string
	v4         *signerV4.Signer
}

func NewAwsSigner(awsFilename, awsProfile, awsRegion, awsService string) (signer *Signer, err error) {
	if err = validateAwsSDKSigner(awsRegion, awsService); err != nil {
		return
	}

	var sess *session.Session
	sess, err = session.NewSession()
	if err != nil {
		return
	}

	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{
				Filename: awsFilename,
				Profile:  awsProfile,
			},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(sess),
			},
		})

	signer = &Signer{
		awsRegion:  awsRegion,
		awsService: awsService,
		v4:         signerV4.NewSigner(creds),
	}

	return
}

func (s *Signer) Sign(req *http.Request, bodyReader io.ReadSeeker, currentTime time.Time) (err error) {
	if s != nil && s.v4 != nil {
		_, err = s.v4.Sign(req, bodyReader, s.awsService, s.awsRegion, time.Now())
		if err != nil {
			return
		}
	} else {
		awsauth.Sign(req)
	}

	return
}

func validateAwsSDKSigner(awsRegion, awsService string) error {
	if awsRegion == "" {
		return errors.New("No AWS region was provided. Cannot sign request.")
	}

	if awsService == "" {
		return errors.New("No AWS service was provided. Cannot sign request.")
	}

	return nil
}

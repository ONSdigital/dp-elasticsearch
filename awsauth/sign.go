package awsauth

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	signerV4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	awsauth "github.com/smartystreets/go-aws-auth"
)

type SignerConfig struct {
	awsProfile   string
	awsRegion    string
	awsSDKSigner bool
	awsService   string
}

func NewSigner(awsSDKSigner bool, awsProfile, awsRegion, awsService string) *SignerConfig {
	return &SignerConfig{
		awsProfile:   awsProfile,
		awsRegion:    awsRegion,
		awsSDKSigner: awsSDKSigner,
		awsService:   awsService,
	}
}

func (sCfg *SignerConfig) Sign(req *http.Request, bodyReader io.ReadSeeker, currentTime time.Time) error {
	if sCfg.awsSDKSigner {
		if err := sCfg.validateAwsSDKSigner(); err != nil {
			return err
		}

		// Remove - attempt to use v2 of go aws sdk
		// cfg, err := config.LoadDefaultConfig(context.Background())
		// if err != nil {
		// 	return err
		// }

		// creds, err := cfg.Credentials.Retrieve(context.Background())
		// if err != nil {
		// 	return err
		// }

		// payloadHash, newReader, err := hashPayload(req.Body)
		// if err != nil {
		// 	return err
		// }
		// req.Body = newReader

		// signer := signerV4.NewSigner()

		// err = signer.SignHTTP(context.Background(), creds, req, payloadHash, sCfg.awsService, sCfg.awsRegion, time.Now())
		// if err != nil {
		// 	return err
		// }

		sess, err := session.NewSession()
		if err != nil {
			return err
		}

		creds := credentials.NewChainCredentials(
			[]credentials.Provider{
				&credentials.EnvProvider{},
				&credentials.SharedCredentialsProvider{
					Filename: "",
					Profile:  sCfg.awsProfile,
				},
				&ec2rolecreds.EC2RoleProvider{
					Client: ec2metadata.New(sess),
				},
			})

		v4Signer := signerV4.NewSigner(creds)
		v4Signer.Sign(req, bodyReader, sCfg.awsService, sCfg.awsRegion, time.Now())
	} else {
		awsauth.Sign(req)
	}

	return nil
}

// TODO - remove hashPayload function used for aws sdk version 2 only
func hashPayload(r io.ReadCloser) (payloadHash string, newReader io.ReadCloser, err error) {
	var payload []byte
	if r == nil {
		payload = []byte("")
	} else {
		payload, err = ioutil.ReadAll(r)
		if err != nil {
			return
		}
		newReader = ioutil.NopCloser(bytes.NewReader(payload))
	}
	hash := sha256.Sum256(payload)
	payloadHash = hex.EncodeToString(hash[:])
	return
}

func (sCfg *SignerConfig) validateAwsSDKSigner() error {
	if sCfg.awsRegion == "" {
		return errors.New("No AWS region was provided. Cannot sign request.")
	}

	if sCfg.awsService == "" {
		return errors.New("No AWS service was provided. Cannot sign request.")
	}

	return nil
}

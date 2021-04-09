package awsauth

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

const (
	envAccessKey       = "AWS_ACCESS_KEY"
	envAccessKeyID     = "AWS_ACCESS_KEY_ID"
	envSecretKey       = "AWS_SECRET_KEY"
	envSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	envSessionToken    = "AWS_SESSION_TOKEN"
)

type Credentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string `json:"Token"`
	Expiration      time.Time
}

// Provider needs to satisfy credentials.Provider interface
type Provider struct {
	credentials.Provider

	expiration time.Time
}

func retrieveCredentials() *credentials.Credentials {
	provider := getProvider()

	return credentials.NewCredentials(provider)
}

func getProvider() credentials.Provider {
	var provider Provider

	return provider
}

// Retrieve produces a set of credentials based on the environment, firstly attempting to retrieve
// credentials from environment variables then attempting to retrieve from IAM Role
func (p Provider) Retrieve() (credentials.Value, error) {
	credValues := credentials.Value{}

	// First use credentials from environment variables
	credValues.AccessKeyID = os.Getenv(envAccessKeyID)
	if credValues.AccessKeyID == "" {
		credValues.AccessKeyID = os.Getenv(envAccessKey)
	}

	credValues.SecretAccessKey = os.Getenv(envSecretAccessKey)
	if credValues.SecretAccessKey == "" {
		credValues.SecretAccessKey = os.Getenv(envSecretKey)
	}

	credValues.SessionToken = os.Getenv(envSessionToken)

	// If there is no Access Key and you are on EC2, get the key from the role
	if (credValues.AccessKeyID == "" || credValues.SecretAccessKey == "") && onEC2() {
		c := *getIAMRoleCredentials()

		credValues.AccessKeyID = c.AccessKeyID
		credValues.SecretAccessKey = c.SecretAccessKey
		credValues.SessionToken = c.SessionToken
		p.expiration = c.Expiration
	}

	return credValues, nil
}

// IsExpired checks to see if the temporary credentials from an IAM role are
// within 4 minutes of expiration (The IAM documentation says that new keys
// will be provisioned 5 minutes before the old keys expire). Credentials
// that do not have an Expiration cannot expire.
func (p Provider) IsExpired() bool {
	if p.expiration.IsZero() {
		// Credentials with no expiration can't expire
		return false
	}

	expireTime := p.expiration.Add(-4 * time.Minute)
	// if t - 4 mins is before now, true
	if expireTime.Before(time.Now()) {
		return true
	} else {
		return false
	}
}

type location struct {
	ec2     bool
	checked bool
}

var loc *location

// onEC2 checks to see if the program is running on an EC2 instance.
// It does this by looking for the EC2 metadata service.
// This caches that information in a struct so that it doesn't waste time.
func onEC2() bool {
	if loc == nil {
		loc = &location{}
	}

	if !(loc.checked) {
		c, err := net.DialTimeout("tcp", "169.254.169.254:80", time.Millisecond*100)
		if err != nil {
			loc.ec2 = false
		} else {
			c.Close()
			loc.ec2 = true
		}

		loc.checked = true
	}

	return loc.ec2
}

func getIAMRoleCredentials() *Credentials {

	roles := getIAMRoleList()

	if len(roles) < 1 {
		return &Credentials{}
	}

	// Use the first role in the list
	role := roles[0]

	url := "http://169.254.169.254/latest/meta-data/iam/security-credentials/"

	// Create the full URL of the role
	var buffer bytes.Buffer
	buffer.WriteString(url)
	buffer.WriteString(role)
	roleURL := buffer.String()

	// Get the role
	roleRequest, err := http.NewRequest("GET", roleURL, nil)
	if err != nil {
		return &Credentials{}
	}

	client := &http.Client{}
	roleResponse, err := client.Do(roleRequest)
	if err != nil {
		return &Credentials{}
	}
	defer roleResponse.Body.Close()

	roleBuffer := new(bytes.Buffer)
	roleBuffer.ReadFrom(roleResponse.Body)

	newCredentials := Credentials{}

	err = json.Unmarshal(roleBuffer.Bytes(), &newCredentials)
	if err != nil {
		return &Credentials{}
	}

	return &newCredentials
}

// getIAMRoleList gets a list of the roles that are available to this instance
func getIAMRoleList() []string {

	var roles []string
	url := "http://169.254.169.254/latest/meta-data/iam/security-credentials/"

	client := &http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return roles
	}

	response, err := client.Do(request)
	if err != nil {
		return roles
	}
	defer response.Body.Close()

	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		roles = append(roles, scanner.Text())
	}

	return roles
}

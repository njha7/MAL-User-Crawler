package main

import (
	"github.com/aws/aws-sdk-go/aws/defaults"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	region = "us-east-1"
)

func main() {
	for {

	}
}

func getSession() *session.Session {
	config := aws.NewConfig().
		WithMaxRetries(3).
		WithRegion(region)

	session := session.Must(session.NewSession(
		config.
			WithCredentials(defaults.CredChain(config, defaults.Handlers())),
	))
	return session
}

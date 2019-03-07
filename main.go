package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/njha7/malutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	serviceName = "MalUserCrawler"
	region      = "us-east-1"
)

func main() {
	session := getSession()
	metricsClient := cloudwatch.New(session)
	for {
		users, err := malutil.GetUsers()
		userCount := 0
		errCount := 0
		if err != nil {
			errCount++
		} else {
			userCount += len(users)
		}
		// TODO enque

		// Publish metrics
		_, err = metricsClient.PutMetricData(&cloudwatch.PutMetricDataInput{
			Namespace: aws.String(serviceName),
			MetricData: []*cloudwatch.MetricDatum{
				// Error count
				&cloudwatch.MetricDatum{
					MetricName: aws.String("GetUserError"),
					Unit:       aws.String(cloudwatch.StandardUnitCount),
					Value:      aws.Float64(float64(errCount)),
				},
				// Discovered user count
				&cloudwatch.MetricDatum{
					MetricName: aws.String("GetUserCount"),
					Unit:       aws.String(cloudwatch.StandardUnitCount),
					Value:      aws.Float64(float64(userCount)),
				},
			},
		})
		if err != nil {
			fmt.Println(err)
		}
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

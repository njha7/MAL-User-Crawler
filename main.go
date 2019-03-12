package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go/service/sqs"

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
	sqsClient := sqs.New(session)
	queueURL := os.Getenv("QUEUE")
	if queueURL == "" {
		panic("Environment variable QUEUE must not be null")
	}
	for {
		users, err := malutil.GetUsers()
		userCount := 0
		errCount := 0
		unqueuedCount := 0
		if err != nil {
			errCount++
		} else {
			userCount += len(users)
		}
		queueMessages := buildQueueMessages(users)
		for _, batch := range queueMessages {
			response, err := sqsClient.SendMessageBatch(&sqs.SendMessageBatchInput{
				QueueUrl: aws.String(queueURL),
				Entries:  batch,
			})
			if err != nil {
				fmt.Println(err)

			} else {
				unqueuedCount += len(response.Failed)
			}
		}

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
				// Unqueued user count
				&cloudwatch.MetricDatum{
					MetricName: aws.String("PutUserError"),
					Unit:       aws.String(cloudwatch.StandardUnitCount),
					Value:      aws.Float64(float64(unqueuedCount)),
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

func buildQueueMessages(users []string) [][]*sqs.SendMessageBatchRequestEntry {
	messages := make([][]*sqs.SendMessageBatchRequestEntry, 0, len(users))
	batch := make([]*sqs.SendMessageBatchRequestEntry, 0, 10)
	for _, user := range users {
		// 10 is the max size of an SQS batch put
		if len(batch)+1 <= 10 {
			batch = append(batch, buildBatchRequestEntry(user))
		} else {
			messages = append(messages, batch)
			batch := make([]*sqs.SendMessageBatchRequestEntry, 0, 10)
			batch = append(batch, buildBatchRequestEntry(user))
		}
	}
	if len(batch) > 0 {
		messages = append(messages, batch)
	}
	return messages
}

func buildBatchRequestEntry(user string) *sqs.SendMessageBatchRequestEntry {
	UUID := uuid.New()
	id := strconv.FormatUint(binary.BigEndian.Uint64(UUID[:]), 16)
	return &sqs.SendMessageBatchRequestEntry{
		Id:          aws.String(string(id)),
		MessageBody: aws.String(user),
	}
}

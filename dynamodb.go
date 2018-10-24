package main

import (
	"log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Worker struct {
	FirstName string `json:first_name`
	LastName string `json:last_name`
	Email string `json:email`
	Address string `json:address`
}

sess, err := session.NewSession(&aws.Config{
	Region: aws.String("us-east-1")},
)

// Create DynamoDB client
svc := dynamodb.New(sess)

func createWorker(*worker Worker) {
	av, err := dynamodbattribute.MarshalMap(*worker)

	input := &dynamodb.PutItemInput{
    Item: av,
    TableName: aws.String("Worker"),
	}

	resp, err = svc.PutItem(input)
	log.Println(resp)
	if err != nil {
    log.Println("Got error calling PutItem:")
    log.Println(err.Error())
    os.Exit(1)
}

log.Println("Successfully added Worker")
}
package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/db"
)

func handleRequest() (events.APIGatewayProxyResponse, error) {

	_, err := db.Init()

	if err != nil {
		panic(err)
	}

	return events.APIGatewayProxyResponse{
		Body:       "OK",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}

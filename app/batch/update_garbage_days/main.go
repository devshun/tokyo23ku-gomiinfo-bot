package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest() (events.APIGatewayProxyResponse, error) {

	fmt.Println("start HandleRequest")

	return events.APIGatewayProxyResponse{}, nil
}

func main() {
	lambda.Start(handleRequest)
}

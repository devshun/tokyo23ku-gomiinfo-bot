package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	db "github.com/devshun/tokyo23ku-gomiinfo-bot/infrastructure"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/infrastructure/mysql"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/usecase"
)

type RequestBody struct {
	Name string
}

func getGarbageDayInfo(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db, err := db.Init()

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	fmt.Println(req)

	var content RequestBody

	err = json.Unmarshal([]byte(req.Body), &content)

	if err != nil {
		panic(err)
	}

	parts := strings.SplitN(content.Name, "区", 2)

	wardName := parts[0] + "区"
	regionName := parts[1]

	m := mysql.NewGarbageDayRepository(db)
	gu := usecase.NewGarbageDayUsecase(m)

	garbageDayInfo, err := gu.GetByAreaNames(wardName, regionName)

	if err != nil {
		panic(err)
	}

	res, err := json.Marshal(garbageDayInfo)

	if err != nil {
		panic(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(res),
	}, nil

}

func main() {
	lambda.Start(getGarbageDayInfo)
}

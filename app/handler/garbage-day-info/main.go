package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	db "github.com/devshun/tokyo23ku-gomiinfo-bot/infrastructure"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/infrastructure/mysql"
)

// TODO: 地域の情報を受け取って以下を返す。
// {
//   "燃えるゴミ": "月曜日、木曜日",
//   "燃えないゴミ": "水曜日",
//   "資源ゴミ": "第1・第3火曜日",
// }
//

type RequestBody struct {
	Name string
}

func getGarbageDayInfo(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	db, err := db.Init()

	if err != nil {
		panic(err)
	}

	fmt.Println(request)

	var content RequestBody

	err = json.Unmarshal([]byte(request.Body), &content)

	if err != nil {
		panic(err)
	}

	parts := strings.SplitN(content.Name, "区", 2)

	wardName := parts[0] + "区"
	regionName := parts[1]

	m := mysql.NewGarbageDayRepository(db)

	garbageDays, err := m.GetByAreaNames(wardName, regionName)

	fmt.Println(garbageDays)

	if err != nil {
		panic(err)
	}

	garbageDaysJSON, err := json.Marshal(garbageDays)

	fmt.Println(string(garbageDaysJSON))

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(garbageDaysJSON),
	}, nil

}

func main() {
	lambda.Start(getGarbageDayInfo)
}

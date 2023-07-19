package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/db"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/domain/model"
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

func getGarbageDayInfo(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	// Add "区" back to the ward name
	wardName := parts[0] + "区"
	regionName := parts[1]

	var garbageDays []model.GarbageDay

	err = db.Preload("Region").Preload("Region.Ward").
		Joins("JOIN regions ON garbage_days.region_id = regions.id").
		Joins("JOIN wards ON regions.ward_id = wards.id").
		Where("wards.name = ? AND regions.name = ?", wardName, regionName).
		Order("garbage_days.garbage_type, garbage_days.day_of_week, garbage_days.week_number_of_month").
		Find(&garbageDays).Error

	garbageDaysJSON, err := json.Marshal(garbageDays)

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

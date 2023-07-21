package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	db "github.com/devshun/tokyo23ku-gomiinfo-bot/infrastructure"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/infrastructure/mysql"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/usecase"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type RequestBody struct {
	Events []struct {
		ReplyToken string `json:"replyToken,omitempty"`
		Type       string `json:"type,omitempty"`
		Timestamp  int64  `json:"timestamp,omitempty"`
		Source     struct {
			Type   string `json:"type,omitempty"`
			UserID string `json:"userId,omitempty"`
		} `json:"source,omitempty"`
		Message struct {
			ID        string  `json:"id,omitempty"`
			Type      string  `json:"type,omitempty"`
			Title     string  `json:"title,omitempty"`
			Address   string  `json:"address,omitempty"`
			Latitude  float64 `json:"latitude,omitempty"`
			Longitude float64 `json:"longitude,omitempty"`
		} `json:"message,omitempty"`
	} `json:"events,omitempty"`
}

func postLineMessage(userid string, message string) error {
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_TOKEN"))
	if err != nil {
		return err
	}

	_, err = bot.PushMessage(userid, linebot.NewTextMessage(message)).Do()

	if err != nil {
		return err
	}

	return nil
}

func getAreaStr(address string) (string, string) {
	parts := strings.Split(address, " ")
	addressParts := strings.Split(parts[1], " ")
	a := addressParts[0]

	// NOTE: とりあえず、○丁目の部分は削除とする
	pattern := `(?P<ward>[^都]+区)(?P<region>[^\d]*[^０-９]*)丁目`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(a)

	paramsMap := make(map[string]string)

	for i, name := range re.SubexpNames() {
		if i > 0 && i <= len(matches) {
			paramsMap[name] = matches[i]
		}
	}

	patternDigits := `[０-９]+$`
	reDigits := regexp.MustCompile(patternDigits)
	region := reDigits.ReplaceAllString(paramsMap["region"], "")

	fmt.Println("ward: ", paramsMap["ward"])
	fmt.Println("region: ", region)

	return paramsMap["ward"], region
}

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	db, err := db.Init()

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	m := mysql.NewGarbageDayRepository(db)
	gu := usecase.NewGarbageDayUsecase(m)

	var event RequestBody

	err = json.Unmarshal([]byte(req.Body), &event)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	// 最初の要素のみ取得
	r := event.Events[0]

	if r.Message.Type == "location" {

		ward, region := getAreaStr(r.Message.Address)

		garbageDayInfo, err := gu.GetByAreaNames(ward, region)

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       err.Error(),
			}, nil
		}

		m := fmt.Sprintf("燃えるゴミ: %s\n燃えないごみ: %s\n資源ごみ: %s",
			garbageDayInfo.Burnable, garbageDayInfo.NonBurnable, garbageDayInfo.Recyclable)

		err = postLineMessage(r.Source.UserID, m)

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       err.Error(),
			}, nil
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}

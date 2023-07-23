package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	db "github.com/devshun/tokyo23ku-gomiinfo-bot/infrastructure"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/infrastructure/mysql"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/usecase"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"golang.org/x/text/width"
	_ "golang.org/x/text/width"
)

// Massaging api request body
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

func postLineMessage(userId string, message string) error {
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_TOKEN"))

	if err != nil {
		return err
	}

	_, err = bot.PushMessage(userId, linebot.NewTextMessage(message)).Do()

	if err != nil {
		return err
	}

	return nil
}

// NOTE: 以下のような文字列から区、地域、丁目を取得する
// 日本、〒130-0021 東京都〇〇区〇〇丁目1−1 〇〇
func getAreaStr(address string) (string, string, int) {

	r, _ := regexp.Compile("東京都(.*?)区(.*?)([０-９]+)")

	match := r.FindStringSubmatch(address)

	var name string
	var region string
	var blockNumber int

	if len(match) > 0 {
		name = match[1] + "区"
	}

	if len(match) > 1 {
		region = match[2]
	}

	if len(match) > 2 {

		// 全角->半角に変換
		s := width.Fold.String(match[3])

		i, err := strconv.Atoi(s)

		if err != nil {
			fmt.Println(err)
		}

		blockNumber = i
	}

	return name, region, blockNumber
}

func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	db, err := db.Init()

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}

	m := mysql.NewGarbageDayRepository(db)
	gu := usecase.NewGarbageDayUsecase(m)

	var event RequestBody

	err = json.Unmarshal([]byte(req.Body), &event)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, err
	}

	// 最初の要素のみ取得
	r := event.Events[0]

	if r.Message.Type == "location" {

		ward, region, blockNumber := getAreaStr(r.Message.Address)

		garbageDayInfo, err := gu.GetByAreaInfo(ward, region, blockNumber)

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       err.Error(),
			}, err
		}

		err = postLineMessage(r.Source.UserID, garbageDayInfo.Format())

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       err.Error(),
			}, err
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

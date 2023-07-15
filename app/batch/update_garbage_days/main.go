package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/devshun/tokyo23ku-gomiinfo-bot.git/db"
)

func updateGarbageDays() error {

	fmt.Println("ゴミ情報の更新を開始します")

	db, err := db.Init()

	if err != nil {
		return err
	}

	defer db.Close()

	fmt.Println("ゴミ情報の更新を終了します")

	return nil
}

func main() {
	lambda.Start(updateGarbageDays)
}

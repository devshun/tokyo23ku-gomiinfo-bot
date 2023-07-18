package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/db"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/domain/model"
)

func getWards() []model.Ward {

	db, err := db.Init()

	if err != nil {
		panic(err)
	}

	var wards []model.Ward

	db.Find(&wards)

	return wards

}

func main() {
	lambda.Start(getWards)
}

package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/batch/import_csv_to_db/importers"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/domain/model"
	db "github.com/devshun/tokyo23ku-gomiinfo-bot/infrastructure"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"gorm.io/gorm"
)

type Config struct {
	name          string
	url           string
	importPackage func(db *gorm.DB, ward model.Ward, records [][]string) error
}

var config = []Config{
	{
		name:          "墨田区",
		url:           "https://www.city.sumida.lg.jp/kuseijoho/sumida_info/opendata/opendata_ichiran/gomirecycle_data/bunbetu_data.files/bunbetu_20151029.csv",
		importPackage: importers.ImportSumidakuCSV,
	},
}

func importCSVToDB() error {

	fmt.Println("INFO: ゴミ情報のインポートを開始します")

	db, err := db.Init()

	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	wg.Add(len(config))

	for _, c := range config {

		go func(c Config) {
			defer wg.Done()

			resp, err := http.Get(c.url)

			if err != nil {
				fmt.Println(err)
			}

			defer resp.Body.Close()

			// ShiftJISをUTF-8に変換
			reader := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())

			r := csv.NewReader(reader)

			records, err := r.ReadAll()

			if err != nil {
				fmt.Println(err)
			}

			var ward model.Ward

			err = db.FirstOrCreate(&ward, model.Ward{Name: c.name}).Error

			if err != nil {
				fmt.Println(err)
			}

			err = c.importPackage(db, ward, records)

			if err != nil {
				fmt.Println(err)
			}

		}(c)
	}

	wg.Wait()

	fmt.Println("INFO: ゴミ情報のインポートを終了します")

	return nil
}

func main() {
	lambda.Start(importCSVToDB)
}

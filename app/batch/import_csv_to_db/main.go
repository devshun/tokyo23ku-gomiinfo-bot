package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/domain/model"
	db "github.com/devshun/tokyo23ku-gomiinfo-bot/infrastructure"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func FindWeekday(s string) (model.Weekday, int, error) {
	for k, v := range model.WeekdayMap {
		// 曜日を取得
		if strings.Contains(s, v) {
			// 第何週目かを取得
			re := regexp.MustCompile(`第(\d)`)

			match := re.FindStringSubmatch(s)

			if len(match) > 0 {
				weekNum, err := strconv.Atoi(match[1])

				if err != nil {
					return 0, 0, err
				}

				return model.Weekday(k), weekNum, nil
			}

			return model.Weekday(k), 0, nil
		}
	}
	return 0, 0, fmt.Errorf("invalid: %s", s)
}

func importCSVToDB() error {

	fmt.Println("INFO: ゴミ情報の更新を開始します")

	db, err := db.Init()

	if err != nil {
		return err
	}

	url := "https://www.city.sumida.lg.jp/kuseijoho/sumida_info/opendata/opendata_ichiran/gomirecycle_data/bunbetu_data.files/bunbetu_20151029.csv"

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// ShiftJISをUTF-8に変換
	reader := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())

	r := csv.NewReader(reader)

	rows, err := r.ReadAll()

	if err != nil {
		return err
	}

	var ward model.Ward

	db.FirstOrCreate(&ward, model.Ward{Name: "墨田区"})

	header := rows[0][1:]

	records := rows[1:]

	for _, record := range records {

		var region model.Region

		db.FirstOrCreate(&region, model.Region{Name: record[0], WardID: ward.ID})

		for i, v := range record[1:] {
			var garbageDay model.GarbageDay

			weekday, weekNum, err := FindWeekday(v)

			if err != nil {
				return err
			}

			garbageType := func() model.GarbageType {
				if header[i] == "燃やすごみの収集曜日" {
					return model.Burnable
				}
				if header[i] == "燃やさないごみの収集曜日" {
					return model.NonBurnable
				}
				if header[i] == "資源物の収集曜日" {
					return model.Recyclable
				}
				return 0
			}()

			db.FirstOrCreate(&garbageDay, model.GarbageDay{RegionID: region.ID, GarbageType: garbageType, DayOfWeek: weekday, WeekNumberOfMonth: weekNum})
		}
	}

	fmt.Println("INFO: ゴミ情報の更新を終了します")

	return nil
}

func main() {
	lambda.Start(importCSVToDB)
}

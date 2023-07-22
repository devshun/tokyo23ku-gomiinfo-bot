package taitoku

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/devshun/tokyo23ku-gomiinfo-bot/domain/model"
	"gorm.io/gorm"
)

var weekdayLabelMap = map[string]model.Weekday{
	"月曜": model.Monday,
	"火曜": model.Tuesday,
	"水曜": model.Wednesday,
	"木曜": model.Thursday,
	"金曜": model.Friday,
	"土曜": model.Saturday,
	"日曜": model.Sunday,
}

var garbageTypeLabelMap = map[string]model.GarbageType{
	"燃やすごみ":   model.Burnable,
	"燃やさないごみ": model.NonBurnable,
	"資源":      model.Recyclable,
}

func getWeekday(s string) (model.Weekday, int, error) {
	for k, v := range weekdayLabelMap {
		// 曜日を取得
		if strings.Contains(s, k) {
			return v, 0, nil
		}
	}
	return 0, 0, fmt.Errorf("invalid: %s", s)
}

func getGarbageType(s string) model.GarbageType {

	v, ok := garbageTypeLabelMap[s]

	if !ok {
		return 0
	}

	return v
}

func ImportTaitokuCSV(db *gorm.DB, ward model.Ward, records [][]string) error {

	header := records[0][1:]

	rows := records[1:]

	var garbageDays []model.GarbageDay

	for _, row := range rows {

		// 浅草1丁目・浅草2丁目 "・"で分割
		names := strings.Split(row[1], "・")

		// 空のRegionのスライスを作成
		var regions []model.Region

		// 各要素をRegionに変換
		for _, name := range names {
			regions = append(regions, model.Region{Name: name, WardID: ward.ID})
		}

		// INSERT REGION
		err := db.Create(&regions).Error

		if err != nil {
			return err
		}

		// 資源、燃えるゴミ
		for i, v := range row[2:4] {

			gerbageType := getGarbageType(header[i+1])

			n := strings.Split(v, "・")

			for _, name := range n {

				weekday, weekNum, err := getWeekday(name)

				if err != nil {
					return err
				}

				for _, r := range regions {
					garbageDays = append(garbageDays, model.GarbageDay{
						RegionID:          r.ID,
						GarbageType:       gerbageType,
						DayOfWeek:         weekday,
						WeekNumberOfMonth: weekNum,
					})
				}

			}
		}

		// 燃えないごみ
		// NOTE: 文字列が特殊なため個別で扱っている。。。
		weekday, err := func() (model.Weekday, error) {
			for k, v := range weekdayLabelMap {
				if strings.Contains(row[4], k) {
					return v, nil
				}
			}
			return 0, fmt.Errorf("invalid: %s", row[4])
		}()

		if err != nil {
			return nil
		}

		weekNums, err := func() ([]int, error) {
			re := regexp.MustCompile("[0-9]+")
			matches := re.FindAllString(row[4], -1)

			var nums []int

			for _, match := range matches {
				num, err := strconv.Atoi(match)
				if err != nil {
					return nil, err
				}
				nums = append(nums, num)
			}
			return nums, nil
		}()

		if err != nil {
			return err
		}

		for _, r := range regions {
			for _, num := range weekNums {
				garbageDays = append(garbageDays, model.GarbageDay{
					RegionID:          r.ID,
					GarbageType:       model.NonBurnable,
					DayOfWeek:         weekday,
					WeekNumberOfMonth: num,
				})
			}
		}
	}

	err := db.Create(&garbageDays).Error

	if err != nil {
		return err
	}

	return nil
}

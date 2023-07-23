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

func getWeekday(s string) (model.Weekday, error) {
	for k, v := range weekdayLabelMap {
		// 曜日を取得
		if strings.Contains(s, k) {
			return v, nil
		}
	}

	return 0, fmt.Errorf("invalid: %s", s)
}

func getGarbageType(s string) model.GarbageType {

	v, ok := garbageTypeLabelMap[s]

	if !ok {
		return 0
	}

	return v
}

func parseRegion(s string) (string, int, error) {
	re := regexp.MustCompile(`(.*?)([0-9０-９一二三四五六七八九十]*)丁目?`)
	match := re.FindStringSubmatch(s)

	if len(match) < 2 {
		return s, 0, nil
	}

	name := match[1]

	var blockNumber int

	if len(match) > 2 && match[2] != "" {
		numStr := match[2]
		var num int
		var err error

		num, err = strconv.Atoi(numStr) // アラビア数字を数値に変換
		if err != nil {
			return "", 0, err
		}

		blockNumber = num
	}

	return name, blockNumber, nil
}

// Header: [索引 町丁名 資源 燃やすごみ 燃やさないごみ]
// Row: [あ 秋葉原 水曜 月曜・木曜 その月の1回目・3回目の土曜日]

func ImportTaitokuCSV(db *gorm.DB, ward model.Ward, records [][]string) error { // nolint:funlen, cyclop

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

			name, blockNumber, err := parseRegion(name)

			if err != nil {
				return err
			}
			regions = append(regions, model.Region{Name: name, BlockNumber: blockNumber, WardID: ward.ID})
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

				weekday, err := getWeekday(name)

				if err != nil {
					return err
				}

				for _, r := range regions {
					garbageDays = append(garbageDays, model.GarbageDay{
						RegionID:    r.ID,
						GarbageType: gerbageType,
						DayOfWeek:   weekday,
					})
				}

			}
		}

		// 燃えないごみ
		// NOTE: 文字列が特殊なため個別で扱っている。。。
		weekday, err := getWeekday(row[4])

		if err != nil {
			return err
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

	// gormのCreateメソッドにスライスを渡してBulkInsert
	err := db.Create(&garbageDays).Error

	if err != nil {
		return err
	}

	return nil
}

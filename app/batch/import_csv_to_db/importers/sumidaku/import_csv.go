package sumidaku

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/devshun/tokyo23ku-gomiinfo-bot/domain/model"
	"gorm.io/gorm"
)

var weekdayLabelMap = map[string]model.Weekday{
	"月曜日": model.Monday,
	"火曜日": model.Tuesday,
	"水曜日": model.Wednesday,
	"木曜日": model.Thursday,
	"金曜日": model.Friday,
	"土曜日": model.Saturday,
	"日曜日": model.Sunday,
}

var garbageTypeLabelMap = map[string]model.GarbageType{
	"燃やすごみの収集曜日":   model.Burnable,
	"燃やさないごみの収集曜日": model.NonBurnable,
	"資源物の収集曜日":     model.Recyclable,
}

var kanjiNumberMap = map[string]int{
	"一": 1,
	"二": 2,
	"三": 3,
	"四": 4,
	"五": 5,
	"六": 6,
	"七": 7,
	"八": 8,
	"九": 9,
	"十": 10,
}

func kanjiToNumber(kanji string) (int, error) {
	if num, ok := kanjiNumberMap[kanji]; ok {
		return num, nil
	}

	return 0, fmt.Errorf("invalid kanji number: %s", kanji)
}

func parseRegion(s string) (string, int, error) {
	re := regexp.MustCompile(`(.*?)([一二三四五六七八九十]*)丁目?`)

	match := re.FindStringSubmatch(s)

	if len(match) < 2 {
		return s, 0, nil
	}

	name := match[1]

	var blockNumber int

	if len(match) > 2 {
		numStr := match[2]

		var err error

		n, err := kanjiToNumber(numStr) // 漢数字を数値に変換

		if err != nil {
			return "", 0, err
		}

		blockNumber = n
	}

	return name, blockNumber, nil
}

// Weekdayを取得
func getWeekday(s string) (model.Weekday, int, error) {
	for k, v := range weekdayLabelMap {
		// 曜日を取得
		if strings.Contains(s, k) {
			// 第何週目かを取得
			re := regexp.MustCompile(`第(\d)`)

			match := re.FindStringSubmatch(s)

			if len(match) > 0 {
				weekNum, err := strconv.Atoi(match[1])

				if err != nil {
					return 0, 0, err
				}

				return v, weekNum, nil
			}

			return v, 0, nil
		}
	}

	return 0, 0, fmt.Errorf("invalid: %s", s)
}

// GarbageTypeを取得
func getGarbageType(s string) model.GarbageType {

	v, ok := garbageTypeLabelMap[s]

	if !ok {
		return 0
	}

	return v
}

// Header: [地域 資源物の収集曜日 燃やすごみの収集曜日 燃やすごみの収集曜日 燃やさないごみの収集曜日 燃やさないごみの収集曜日]
// Row: [吾妻橋 土曜日 火曜日 金曜日 第1月曜日 第3月曜日]

func ImportSumidakuCSV(db *gorm.DB, ward model.Ward, records [][]string) error {

	header := records[0][1:]

	rows := records[1:]

	var garbageDays []model.GarbageDay

	for _, row := range rows {

		var region model.Region

		name, blockNumber, err := parseRegion(row[0])

		if err != nil {
			return err
		}

		// INSERT REGION
		err = db.FirstOrCreate(&region, model.Region{Name: name, BlockNumber: blockNumber, WardID: ward.ID}).Error

		if err != nil {
			return err
		}

		for i, v := range row[1:] {

			weekday, weekNum, err := getWeekday(v)

			if err != nil {
				return err
			}

			garbageType := getGarbageType(header[i])

			garbageDays = append(garbageDays, model.GarbageDay{Region: region, GarbageType: garbageType, DayOfWeek: weekday, WeekNumberOfMonth: weekNum})
		}
	}

	// BULK INSERT GARBAGE_DAYS
	err := db.Create(garbageDays).Error

	if err != nil {
		return err
	}

	return nil
}

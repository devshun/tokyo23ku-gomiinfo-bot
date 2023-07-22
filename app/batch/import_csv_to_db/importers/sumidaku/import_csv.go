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

func ImportSumidakuCSV(db *gorm.DB, ward model.Ward, records [][]string) error {

	header := records[0][1:]

	rows := records[1:]

	var garbageDays []model.GarbageDay

	for _, row := range rows {

		var region model.Region

		// INSERT REGION
		err := db.FirstOrCreate(&region, model.Region{Name: row[0], WardID: ward.ID}).Error

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

package importers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/devshun/tokyo23ku-gomiinfo-bot/domain/model"
	"gorm.io/gorm"
)

// 曜日を取得
func getWeekDay(s string) (model.Weekday, int, error) {
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

// GarbageTypeを取得
func getGarbageType(s string) model.GarbageType {

	switch s {
	case "燃やすごみの収集曜日":
		return model.Burnable
	case "燃やさないごみの収集曜日":
		return model.NonBurnable
	case "資源物の収集曜日":
		return model.Recyclable
	default:
		return 0
	}
}

func ImportSumidakuCSV(db *gorm.DB, ward model.Ward, records [][]string) error {

	startTime := time.Now()

	header := records[0][1:]

	rows := records[1:]

	var garbageDays []model.GarbageDay

	for _, row := range rows {

		var region model.Region

		err := db.FirstOrCreate(&region, model.Region{Name: row[0], WardID: ward.ID}).Error

		if err != nil {
			return err
		}

		for i, v := range row[1:] {

			weekday, weekNum, err := getWeekDay(v)

			if err != nil {
				return err
			}

			garbageType := getGarbageType(header[i])

			garbageDays = append(garbageDays, model.GarbageDay{Region: region, GarbageType: garbageType, DayOfWeek: weekday, WeekNumberOfMonth: weekNum})
		}
	}

	// BULK INSERT
	err := db.Create(garbageDays).Error

	if err != nil {
		return err
	}

	elapsedTime := time.Since(startTime)

	fmt.Printf("INFO: 墨田区インポートにかかった時間: %s\n", elapsedTime)

	return nil
}

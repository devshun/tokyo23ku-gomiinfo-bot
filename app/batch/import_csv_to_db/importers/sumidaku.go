package importers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/devshun/tokyo23ku-gomiinfo-bot/domain/model"
	"gorm.io/gorm"
)

func ImportSumidakuCSV(db *gorm.DB, ward model.Ward, records [][]string) error {
	header := records[0][1:]

	rows := records[1:]

	for _, row := range rows {

		var region model.Region

		db.FirstOrCreate(&region, model.Region{Name: row[0], WardID: ward.ID})

		for i, v := range row[1:] {
			var garbageDay model.GarbageDay

			weekday, weekNum, err := FindWeekday(v)

			if err != nil {
				return err
			}

			garbageType := func() model.GarbageType {
				switch header[i] {

				case "燃やすごみの収集曜日":
					return model.Burnable

				case "燃やさないごみの収集曜日":
					return model.NonBurnable

				case "資源物の収集曜日":
					return model.Recyclable

				default:
					return 0
				}
			}()

			db.FirstOrCreate(&garbageDay, model.GarbageDay{RegionID: region.ID, GarbageType: garbageType, DayOfWeek: weekday, WeekNumberOfMonth: weekNum})
		}
	}
	return nil
}

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

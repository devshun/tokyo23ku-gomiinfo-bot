package mysql

import (
	"github.com/devshun/tokyo23ku-gomiinfo-bot/domain/model"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/domain/repository"
	"gorm.io/gorm"
)

type garbageDayRepository struct {
	DB *gorm.DB
}

func NewGarbageDayRepository(db *gorm.DB) repository.GarbageDayRepository {
	return &garbageDayRepository{DB: db}
}

func (gr *garbageDayRepository) GetByAreaInfo(ward string, region string, blockNumber int) ([]model.GarbageDay, error) {
	var garbageDays []model.GarbageDay

	err := gr.DB.Preload("Region").Preload("Region.Ward").
		Joins("JOIN regions ON garbage_days.region_id = regions.id").
		Joins("JOIN wards ON regions.ward_id = wards.id").
		// NOTE: ~ 丁目まで一致するレコード or 地域まで一致するレコードを取得
		// TODO: 丁目まで一致するレコードの方が優先度が高いので、丁目まで一致するレコードを先に取得するようにする
		Where("(wards.name = ? AND regions.name = ? AND regions.block_number = ?) OR (wards.name = ? AND regions.name = ? AND regions.block_number = 0)", ward, region, blockNumber, ward, region).
		Order("garbage_days.garbage_type, garbage_days.day_of_week, garbage_days.week_number_of_month").
		Find(&garbageDays).Error

	if err != nil {
		return nil, err
	}

	return garbageDays, nil
}

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

func (gr *garbageDayRepository) GetByAreaNames(wardName string, regionName string) ([]model.GarbageDay, error) {

	var garbageDays []model.GarbageDay

	err := gr.DB.Preload("Region").Preload("Region.Ward").
		Joins("JOIN regions ON garbage_days.region_id = regions.id").
		Joins("JOIN wards ON regions.ward_id = wards.id").
		Where("wards.name = ? AND regions.name = ?", wardName, regionName).
		Order("garbage_days.garbage_type, garbage_days.day_of_week, garbage_days.week_number_of_month").
		Find(&garbageDays).Error

	if err != nil {
		return nil, err
	}

	return garbageDays, nil
}

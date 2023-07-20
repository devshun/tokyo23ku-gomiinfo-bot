package repository

import (
	"github.com/devshun/tokyo23ku-gomiinfo-bot/domain/model"
)

type GarbageDayRepository interface {
	GetByAreaNames(wardName string, regionName string) ([]model.GarbageDay, error)
}

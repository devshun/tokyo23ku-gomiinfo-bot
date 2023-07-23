package repository

import (
	"github.com/devshun/tokyo23ku-gomiinfo-bot/domain/model"
)

type GarbageDayRepository interface {
	GetByAreaInfo(wardName string, regionName string, blockNumber int) ([]model.GarbageDay, error)
}

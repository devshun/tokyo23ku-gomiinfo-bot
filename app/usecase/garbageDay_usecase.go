package usecase

import (
	"github.com/devshun/tokyo23ku-gomiinfo-bot/domain/model"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/domain/repository"
)

type garbageDayUsecase struct {
	garbageDayRepo repository.GarbageDayRepository
}

func NewGarbageDayUsecase(garbageDayRepository repository.GarbageDayRepository) *garbageDayUsecase {
	return &garbageDayUsecase{
		garbageDayRepo: garbageDayRepository,
	}
}

func (gu *garbageDayUsecase) GetByAreaNames(wardName string, regionName string) ([]model.GarbageDay, error) {
	return gu.garbageDayRepo.GetByAreaNames(wardName, regionName)
}

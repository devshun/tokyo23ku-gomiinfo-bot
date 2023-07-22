package usecase

import (
	"fmt"
	"strings"

	"github.com/devshun/tokyo23ku-gomiinfo-bot/domain/model"
	"github.com/devshun/tokyo23ku-gomiinfo-bot/domain/repository"
)

type garbageDayUsecase struct {
	garbageDayRepo repository.GarbageDayRepository
}

type GarbageDayInfo struct {
	Burnable    string `json:"燃えるゴミ,omitempty"`
	NonBurnable string `json:"燃えないゴミ,omitempty"`
	Recyclable  string `json:"資源ゴミ,omitempty"`
}

func NewGarbageDayUsecase(garbageDayRepository repository.GarbageDayRepository) *garbageDayUsecase {
	return &garbageDayUsecase{
		garbageDayRepo: garbageDayRepository,
	}
}

func (gu *garbageDayUsecase) GetByAreaNames(wardName string, regionName string) (GarbageDayInfo, error) {
	garbageDays, err := gu.garbageDayRepo.GetByAreaNames(wardName, regionName)

	if err != nil {
		return GarbageDayInfo{}, err
	}

	burnableDays := []string{}
	nonBurnableDays := []string{}
	recyclableDays := []string{}

	for _, garbageDay := range garbageDays {

		dayOfWeek := garbageDay.Format()

		switch garbageDay.GarbageType {

		case model.Burnable:
			burnableDays = append(burnableDays, dayOfWeek)
		case model.NonBurnable:
			nonBurnableDays = append(nonBurnableDays, dayOfWeek)
		case model.Recyclable:
			recyclableDays = append(recyclableDays, dayOfWeek)
		}
	}

	return GarbageDayInfo{
		Burnable:    strings.Join(burnableDays, "、"),
		NonBurnable: strings.Join(nonBurnableDays, "、"),
		Recyclable:  strings.Join(recyclableDays, "、"),
	}, nil
}

func (gi *GarbageDayInfo) Format() string {

	if len(gi.Burnable) == 0 && len(gi.NonBurnable) == 0 && len(gi.Recyclable) == 0 {
		return "ごみ収集日のデータが見つかりません"
	}

	m := fmt.Sprintf("燃えるゴミ: %s\n燃えないごみ: %s\n資源ごみ: %s",
		gi.Burnable, gi.NonBurnable, gi.Recyclable)

	return m
}

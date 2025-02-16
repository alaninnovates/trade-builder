package loaders

import (
	"github.com/disgoorg/json"
	"os"
)

var beequipDataFile = "assets/data/beequipData.json"

type BeequipData struct {
	Buffs   []string
	Debuffs []string
	Ability []string
	Bonuses []string
}

var beequipData map[string]BeequipData

func LoadBeequips() {
	var beequips map[string]BeequipData
	data, err := os.ReadFile(beequipDataFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &beequips)
	if err != nil {
		panic(err)
	}
	beequipData = beequips
}

func GetBeequip(beequipName string) BeequipData {
	if beequipData == nil {
		LoadBeequips()
	}
	return beequipData[beequipName]
}

func GetBeequipBuffs(beequipName string) []string {
	return GetBeequip(beequipName).Buffs
}

func GetBeequipDebuffs(beequipName string) []string {
	return GetBeequip(beequipName).Debuffs
}

func GetBeequipAbility(beequipName string) []string {
	return GetBeequip(beequipName).Ability
}

func GetBeequipBonuses(beequipName string) []string {
	return GetBeequip(beequipName).Bonuses
}

func GetAllBeequips() []string {
	if beequipData == nil {
		LoadBeequips()
	}
	var beequips []string
	for beequip := range beequipData {
		beequips = append(beequips, beequip)
	}
	return beequips
}

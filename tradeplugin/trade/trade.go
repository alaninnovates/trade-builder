package trade

import (
	"go.mongodb.org/mongo-driver/bson"
)

type BeequipInProgressStep string

const (
	BeequipInProgressStepNone    BeequipInProgressStep = ""
	BeequipInProgressStepBuffs   BeequipInProgressStep = "buffs"
	BeequipInProgressStepDebuffs BeequipInProgressStep = "debuffs"
	BeequipInProgressStepAbility BeequipInProgressStep = "ability"
	BeequipInProgressStepBonuses BeequipInProgressStep = "bonuses"
	BeequipInProgressStepWaxes   BeequipInProgressStep = "waxes"
)

type Sticker struct {
	Name     string
	Quantity int
}

type Beequip struct {
	Name      string
	Buffs     map[string]int
	Debuffs   map[string]int
	Ability   map[string]bool
	Bonuses   map[string]int
	Potential int
	Waxes     []string
}

type Trade struct {
	lookingFor            []interface{}
	offering              []interface{}
	beequipInProgressStep BeequipInProgressStep
	beequipInProgresType  string
	beequipInProgress     Beequip
}

func NewTrade() *Trade {
	return &Trade{}
}

func (t *Trade) AddLookingForSticker(s Sticker) {
	t.lookingFor = append(t.lookingFor, s)
}

func (t *Trade) AddOfferingSticker(s Sticker) {
	t.offering = append(t.offering, s)
}

func (t *Trade) AddLookingForBeequip(b Beequip) {
	t.lookingFor = append(t.lookingFor, b)
}

func (t *Trade) AddOfferingBeequip(b Beequip) {
	t.offering = append(t.offering, b)
}

func (t *Trade) SetBeequipInProgressType(ty string) {
	t.beequipInProgresType = ty
}

func (t *Trade) GetBeequipInProgressType() string {
	return t.beequipInProgresType
}

func (t *Trade) SetBeequipInProgressStep(step BeequipInProgressStep) {
	t.beequipInProgressStep = step
}

func (t *Trade) IsBeequipInProgress() bool {
	return t.beequipInProgressStep != ""
}

func (t *Trade) GetBeequipInProgressStep() BeequipInProgressStep {
	return t.beequipInProgressStep
}

func (t *Trade) GetBeequipInProgressData() Beequip {
	return t.beequipInProgress
}

func (t *Trade) SetBeequipInProgressData(b Beequip) {
	t.beequipInProgress = b
}

func (t *Trade) Remove(category string, name string) {
	var stickers []interface{}
	if category == "lf" {
		stickers = t.lookingFor
	} else {
		stickers = t.offering
	}
	for i, sticker := range stickers {
		if s, ok := sticker.(Sticker); ok && s.Name == name {
			stickers = append(stickers[:i], stickers[i+1:]...)
			break
		}
	}
	if category == "lf" {
		t.lookingFor = stickers
	} else {
		t.offering = stickers
	}
}

func (t *Trade) GetLookingFor() []interface{} {
	return t.lookingFor
}

func (t *Trade) GetOffering() []interface{} {
	return t.offering
}

func (t *Trade) ToBson() bson.D {
	lookingFor := bson.D{}
	for _, sticker := range t.lookingFor {
		if s, ok := sticker.(Sticker); ok {
			lookingFor = append(lookingFor, bson.E{Key: s.Name, Value: s.Quantity})
		}
	}
	offering := bson.D{}
	for _, sticker := range t.offering {
		if s, ok := sticker.(Sticker); ok {
			offering = append(offering, bson.E{Key: s.Name, Value: s.Quantity})
		}
	}
	return bson.D{{"lookingFor", lookingFor}, {"offering", offering}}
}

func (t *Trade) FromBson(doc bson.D) {
	for _, elem := range doc {
		switch elem.Key {
		case "lookingFor":
			lookingFor := elem.Value.(bson.D)
			for _, sticker := range lookingFor {
				t.lookingFor = append(t.lookingFor, Sticker{Name: sticker.Key, Quantity: int(sticker.Value.(int32))})
			}
		case "offering":
			offering := elem.Value.(bson.D)
			for _, sticker := range offering {
				t.offering = append(t.offering, Sticker{Name: sticker.Key, Quantity: int(sticker.Value.(int32))})
			}
		}
	}
}

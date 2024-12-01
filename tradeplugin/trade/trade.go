package trade

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Sticker struct {
	Name     string
	Quantity int
}

type Trade struct {
	lookingFor []Sticker
	offering   []Sticker
}

func NewTrade() *Trade {
	return &Trade{}
}

func (t *Trade) AddLookingFor(name string, quantity int) {
	t.lookingFor = append(t.lookingFor, Sticker{name, quantity})
}

func (t *Trade) AddOffering(name string, quantity int) {
	t.offering = append(t.offering, Sticker{name, quantity})
}

func (t *Trade) Remove(category string, name string) {
	var stickers []Sticker
	if category == "lf" {
		stickers = t.lookingFor
	} else {
		stickers = t.offering
	}
	for i, sticker := range stickers {
		if sticker.Name == name {
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

func (t *Trade) GetLookingFor() []Sticker {
	return t.lookingFor
}

func (t *Trade) GetOffering() []Sticker {
	return t.offering

}

func (t *Trade) ToBson() bson.D {
	lookingFor := bson.D{}
	for _, sticker := range t.lookingFor {
		lookingFor = append(lookingFor, bson.E{Key: sticker.Name, Value: sticker.Quantity})
	}
	offering := bson.D{}
	for _, sticker := range t.offering {
		offering = append(offering, bson.E{Key: sticker.Name, Value: sticker.Quantity})
	}
	return bson.D{{"lookingFor", lookingFor}, {"offering", offering}}
}

func (t *Trade) FromBson(doc bson.D) {
	for _, elem := range doc {
		switch elem.Key {
		case "lookingFor":
			lookingFor := elem.Value.(bson.D)
			for _, sticker := range lookingFor {
				t.lookingFor = append(t.lookingFor, Sticker{sticker.Key, int(sticker.Value.(int32))})
			}
		case "offering":
			offering := elem.Value.(bson.D)
			for _, sticker := range offering {
				t.offering = append(t.offering, Sticker{sticker.Key, int(sticker.Value.(int32))})
			}
		}
	}
}

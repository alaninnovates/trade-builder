package database

import (
	"encoding/json"
	"os"
)

type CachedSticker struct {
	Id       string
	Quantity int
}

type CachedTradeSide map[int]CachedSticker

type CachedUser struct {
	Id         string
	LookingFor CachedTradeSide
	Offering   CachedTradeSide
}

type JsonCache struct {
}

func NewJsonCache() *JsonCache {
	return &JsonCache{}
}

func (j *JsonCache) SaveTrades(fileName string, cachedUsers []CachedUser) error {
	jsonHive, err := json.MarshalIndent(cachedUsers, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	_, err = f.Write(jsonHive)
	if err != nil {
		return err
	}
	return nil
}

func (j *JsonCache) LoadTrades(fileName string) ([]CachedUser, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	var cachedUsers []CachedUser
	err = json.NewDecoder(f).Decode(&cachedUsers)
	if err != nil {
		return nil, err
	}
	return cachedUsers, nil
}

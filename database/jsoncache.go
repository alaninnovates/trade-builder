package database

import (
	"encoding/json"
	"os"
)

type CachedBee struct {
	Id       string
	Level    int
	Gifted   bool
	Beequip  string
	Mutation string
}

type CachedHive map[int]CachedBee

type CachedUser struct {
	Id   string
	Hive CachedHive
}

type JsonCache struct {
}

func NewJsonCache() *JsonCache {
	return &JsonCache{}
}

func (j *JsonCache) SaveHives(fileName string, cachedUsers []CachedUser) error {
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

func (j *JsonCache) LoadHives(fileName string) ([]CachedUser, error) {
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

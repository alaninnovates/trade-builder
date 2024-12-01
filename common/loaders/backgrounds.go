package loaders

import (
	"github.com/fogleman/gg"
	"golang.org/x/exp/rand"
	"image"
	"os"
)

var cachedBackgrounds = make(map[int]image.Image)

func GetRandomBackground() image.Image {
	if len(cachedBackgrounds) == 0 {
		bgCnt, err := os.ReadDir("assets/backgrounds")
		if err != nil {
			panic(err)
		}
		for i, bg := range bgCnt {
			img, err := gg.LoadImage("assets/backgrounds/" + bg.Name())
			if err != nil {
				panic(err)
			}
			cachedBackgrounds[i] = img
		}
	}
	return cachedBackgrounds[rand.Intn(len(cachedBackgrounds))]
}

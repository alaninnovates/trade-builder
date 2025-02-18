package loaders

import (
	"github.com/fogleman/gg"
	"image"
	"os"
)

var cachedWaxes = make(map[string]image.Image)

func GetWaxImage(waxName string) image.Image {
	if len(cachedWaxes) == 0 {
		files, err := os.ReadDir("assets/waxes")
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			img, err := gg.LoadImage("assets/waxes/" + file.Name())
			if err != nil {
				panic(err)
			}
			cachedWaxes[file.Name()[:len(file.Name())-4]] = img
		}
	}
	return cachedWaxes[waxName]
}

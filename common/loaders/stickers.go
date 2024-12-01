package loaders

import (
	"github.com/fogleman/gg"
	"golang.org/x/image/draw"
	"image"
	"os"
)

var stickerDirectories = []string{"assets/stickers", "assets/hive_skins", "assets/cub_skins", "assets/vouchers"}
var cachedStickers = make(map[string]image.Image)

func resize(im image.Image) *image.NRGBA {
	rgbImage := im.(*image.NRGBA)
	dst := image.NewNRGBA(image.Rect(0, 0, 100, 100))
	draw.NearestNeighbor.Scale(dst, dst.Bounds(), rgbImage, rgbImage.Bounds(), draw.Over, nil)
	return dst
}

func GetStickerImage(stickerName string) image.Image {
	if len(cachedStickers) == 0 {
		for _, dir := range stickerDirectories {
			files, err := os.ReadDir(dir)
			if err != nil {
				panic(err)
			}
			for _, file := range files {
				img, err := gg.LoadImage(dir + "/" + file.Name())
				if err != nil {
					panic(err)
				}
				cachedStickers[file.Name()[:len(file.Name())-4]] = img
			}
		}
	}
	return cachedStickers[stickerName]
}

func GetAllStickers() []string {
	var files []string
	for _, dir := range stickerDirectories {
		filesInDir, err := os.ReadDir(dir)
		if err != nil {
			panic(err)
		}
		for _, file := range filesInDir {
			// remove .png extension
			files = append(files, file.Name()[:len(file.Name())-4])
		}
	}
	return files
}

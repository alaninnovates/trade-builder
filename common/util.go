package common

import (
	"fmt"
	"image"
	"image/png"
	"io"
)

func ImageToPipe(image image.Image) *io.PipeReader {
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		if err := png.Encode(w, image); err != nil {
			fmt.Println(err)
		}
	}()
	return r
}

func ArrayIncludes[T comparable](array []T, value T) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

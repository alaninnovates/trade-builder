package common

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"strconv"
	"strings"
	"time"
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

func ParseHHMM(input string) (time.Duration, error) {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid format")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	duration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
	return duration, nil
}

func ParseMMDD(input string) (time.Duration, error) {
	parts := strings.Split(input, "/")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid format")
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	day, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	now := time.Now().UTC()
	year := now.Year()
	expire := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if expire.Before(now) {
		year++
		expire = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	}
	duration := expire.Sub(now)
	return duration, nil
}

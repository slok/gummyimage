package gummyimage

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

type Gummy struct {
	img *image.RGBA
}

func NewDefaultGummy(w, h int) (*Gummy, error) {
	return NewGummy(0, 0, w, h, randColor(255))
}

func NewGummy(x, y, w, h int, gummyColor color.Color) (*Gummy, error) {

	img, err := createImg(x, y, w, h, gummyColor)

	if err != nil {
		return nil, err
	}

	return &Gummy{
		img: img,
	}, nil
}

func (g *Gummy) SavePng(path string) error {
	file, err := os.Create(path)

	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, g.img)
}

func createImg(x, y, w, h int, gummyColor color.Color) (*image.RGBA, error) {
	img := image.NewRGBA(image.Rect(x, y, w, h))

	// Colorize!
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
			img.Set(x, y, gummyColor)
		}
	}

	return img, nil

}

func randColor(alpha int) color.Color {

	random := func(min, max int) int {
		rand.Seed(time.Now().UnixNano())
		return rand.Intn(max-min) + min
	}

	r := uint8(random(0, 255))
	g := uint8(random(0, 255))
	b := uint8(random(0, 255))

	return color.RGBA{r, g, b, uint8(alpha)}

}

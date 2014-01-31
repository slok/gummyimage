package gummyimage

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
)

var (
	font = loadFont("./gummyimage/DroidSans.ttf")
)

func loadFont(path string) *truetype.Font {
	bs, _ := ioutil.ReadFile(path)
	f, _ := truetype.Parse(bs)
	return f
}

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

// Color in HEX format: FAFAFA
func (g *Gummy) DrawText(text, textColor string, fontSize, xPosition, yPosition int) error {

	fc := freetype.NewContext()
	fc.SetDst(g.img)
	fc.SetFont(font)
	fc.SetClip(g.img.Bounds())

	// Color parsing
	cr, _ := strconv.ParseUint(textColor[:2], 16, 64)
	cg, _ := strconv.ParseUint(textColor[2:4], 16, 64)
	cb, _ := strconv.ParseUint(textColor[4:], 16, 64)
	c := image.NewUniform(color.RGBA{R: uint8(cr), G: uint8(cg), B: uint8(cb), A: 255})

	fc.SetSrc(c)
	fc.SetFontSize(float64(fontSize))

	_, err := fc.DrawString(text, freetype.Pt(xPosition, yPosition))

	return err
}

// Color in HEX format: FAFAFA
func (g *Gummy) DrawTextSize(color string) error {

	text := fmt.Sprintf("%dx%d", g.img.Rect.Max.X, g.img.Rect.Max.Y)
	// I can get the text final size so more or less center the text with this
	// awful stuff :/

	size := 0
	fontFactor := 0

	if g.img.Rect.Max.X < g.img.Rect.Max.Y {
		size = g.img.Rect.Max.X
		fontFactor = g.img.Rect.Max.Y / g.img.Rect.Max.X * 2
	} else {
		size = g.img.Rect.Max.Y
		fontFactor = g.img.Rect.Max.X / g.img.Rect.Max.Y * 2
	}

	textSpace := (size - (size / 10 * 2))
	fontSize := textSpace / len(text) * fontFactor
	x := g.img.Rect.Max.X/2 - textSpace/2 - fontSize/8*(fontFactor+1)
	y := g.img.Rect.Max.Y/2 + textSpace/10 + fontSize/16*(fontFactor+1)

	/*
	   text := fmt.Sprintf("%dx%d", g.img.Rect.Max.X, g.img.Rect.Max.Y)
	   // I can get the text final size so more or less center the text with this
	   // awful stuff :/

	   // Get minimun size for correct font
	   size := g.img.Rect.Max.Y
	   fontFactor := 2

	   if g.img.Rect.Max.X < g.img.Rect.Max.Y {
	       size = g.img.Rect.Max.X
	   }

	   // Get total space for the font
	   textSpace := (size - (size / 10 * 2))

	   // Get the font size based on the size of the text (manual approximation)
	   fontSize := textSpace / len(text) * fontFactor

	   //Get approx the center based on the calculated data
	   x := g.img.Rect.Max.X/2 - textSpace/2 - fontSize/2
	   y := g.img.Rect.Max.Y/2 + textSpace/10 + fontSize/4

	   //x := g.img.Rect.Max.X/2 - ((g.img.Rect.Max.X - (g.img.Rect.Max.X / 10 * 2)) / 2) - fontSize/8
	   //y := g.img.Rect.Max.Y/2 + ((g.img.Rect.Max.Y - (g.img.Rect.Max.Y / 10 * 2)) / 10) + fontSize/8
	*/

	return g.DrawText(
		text,
		color,
		fontSize,
		x,
		y,
	)
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

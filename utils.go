package gummyimage

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
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
	img   *image.RGBA
	color *color.Color
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
		img:   img,
		color: &gummyColor,
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
// If "" the color of the text is black or white depending on the brightness of the bg
func (g *Gummy) DrawTextSize(textColor string) error {

	// Get black or white depending on the background
	if textColor == "" {
		c := (*g.color).(color.RGBA)
		if blackWithBackground(float64(c.R), float64(c.G), float64(c.B)) {
			textColor = "000000"
		} else {
			textColor = "FFFFFF"
		}
	}

	text := fmt.Sprintf("%dx%d", g.img.Rect.Max.X, g.img.Rect.Max.Y)

	// I can't get the text final size so more or less center the text with this
	// manual awful stuff :/
	size := g.img.Rect.Max.Y

	if g.img.Rect.Max.X < g.img.Rect.Max.Y {
		size = g.img.Rect.Max.X
	}

	textSize := (size - (size / 10 * 2))
	fontSize := textSize / len(text) * 2

	x := g.img.Rect.Max.X/2 - textSize/2 - fontSize/8
	y := g.img.Rect.Max.Y/2 + textSize/10 + fontSize/16

	return g.DrawText(
		text,
		textColor,
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

func inverseColor(r, g, b int) (rr, rg, rb int) {
	rr = 255 - r
	rg = 255 - g
	rb = 255 - b

	return
}

// Returns false if white text with that background
// Rrturns true if black text with that background
// Calculates based on the brightness
// Source: http://stackoverflow.com/a/2241471
func blackWithBackground(r, g, b float64) bool {

	perceivedBrightness := func(r, g, b float64) int {
		return int(math.Sqrt(r*r*0.241 + g*g*0.691 + b*b*0.068))
	}

	return perceivedBrightness(r, g, b) > 130
}

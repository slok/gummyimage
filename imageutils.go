package gummyimage

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
)

type Gummy struct {
	Img   *image.RGBA
	Color *color.Color
	Font  *truetype.Font
}

// Color in HEX format: FAFAFA
// If hexColor = "" then random color
func NewDefaultGummy(w, h int, hexColor string) (*Gummy, error) {
	var bgColor color.Color
	if hexColor == "" {
		bgColor = randColor(255)

	} else {
		cr, _ := strconv.ParseUint(string(hexColor[:2]), 16, 64)
		cg, _ := strconv.ParseUint(string(hexColor[2:4]), 16, 64)
		cb, _ := strconv.ParseUint(string(hexColor[4:]), 16, 64)
		bgColor = color.RGBA{R: uint8(cr), G: uint8(cg), B: uint8(cb), A: 255}
	}

	return NewGummy(0, 0, w, h, bgColor)
}

func NewGummy(x, y, w, h int, gummyColor color.Color) (*Gummy, error) {

	img, err := createImg(x, y, w, h, gummyColor)

	if err != nil {
		return nil, err
	}

	return &Gummy{
		Img:   img,
		Color: &gummyColor,
		Font:  nil,
	}, nil
}

/*
Gets the image in the specified format (JPEG, GIF or PNG) in the specified writer
*/
func (g *Gummy) get(format string, r io.Writer) error {

	switch format {
	case "jpeg", "JPEG":
		jpeg.Encode(r, g.Img, nil)
	case "png", "PNG":
		png.Encode(r, g.Img)
	case "gif", "GIF":
		gif.Encode(r, g.Img, nil)
	default:
		return errors.New("Wrong format")
	}

	return nil
}

func (g *Gummy) SaveJpeg(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return g.get("jpeg", file)
}

func (g *Gummy) SaveGif(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return g.get("gif", file)
}

func (g *Gummy) SavePng(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return g.get("png", file)
}

func (g *Gummy) GetJpeg() ([]byte, error) {
	b := new(bytes.Buffer)
	err := g.get("jpeg", b)
	return b.Bytes(), err
}

func (g *Gummy) GetGif() ([]byte, error) {
	b := new(bytes.Buffer)
	err := g.get("gif", b)
	return b.Bytes(), err
}

func (g *Gummy) GetPng() ([]byte, error) {
	b := new(bytes.Buffer)
	err := g.get("png", b)
	return b.Bytes(), err
}

// Color in HEX format: FAFAFA
func (g *Gummy) DrawText(text, textColor string, fontSize, xPosition, yPosition int) error {

	fc := freetype.NewContext()
	fc.SetDst(g.Img)
	fc.SetFont(g.Font)
	fc.SetClip(g.Img.Bounds())

	// Color parsing
	cr, _ := strconv.ParseUint(string(textColor[:2]), 16, 64)
	cg, _ := strconv.ParseUint(string(textColor[2:4]), 16, 64)
	cb, _ := strconv.ParseUint(string(textColor[4:]), 16, 64)
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
		c := (*g.Color).(color.RGBA)
		if blackWithBackground(float64(c.R), float64(c.G), float64(c.B)) {
			textColor = "000000"
		} else {
			textColor = "FFFFFF"
		}
	}

	text := fmt.Sprintf("%dx%d", g.Img.Rect.Max.X, g.Img.Rect.Max.Y)

	// I can't get the text final size so more or less center the text with this
	// manual awful stuff :/
	size := g.Img.Rect.Max.Y

	if g.Img.Rect.Max.X < g.Img.Rect.Max.Y {
		size = g.Img.Rect.Max.X
	}

	textSize := (size - (size / 10 * 2))
	fontSize := textSize / len(text) * 2

	x := g.Img.Rect.Max.X/2 - textSize/2 - fontSize/8
	y := g.Img.Rect.Max.Y/2 + textSize/10 + fontSize/16

	return g.DrawText(
		text,
		textColor,
		fontSize,
		x,
		y,
	)
}

func LoadFont(path string) (*truetype.Font, error) {
	bs, err := ioutil.ReadFile(path)

	// quick debug
	if err != nil {
		fmt.Println(err)
	}
	f, err := truetype.Parse(bs)
	return f, err
}

func (g *Gummy) SetFont(path string) error {
	f, err := LoadFont(path)
	g.Font = f
	return err
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

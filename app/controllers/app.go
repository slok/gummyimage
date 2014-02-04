package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"code.google.com/p/freetype-go/freetype/truetype"
	"github.com/robfig/revel"
	"github.com/slok/gummyimage"
)

type Application struct {
	*revel.Controller
}

type ImageResponse struct {
	sizeX   int
	sizeY   int
	bgColor string
	fgColor string
	text    string
	format  string
}

// Global variable
var (
	font              *truetype.Font
	regularSizeRegex  = regexp.MustCompile(`^(.+)[xX](.+)$`)
	aspectSizeRegex   = regexp.MustCompile(`^(.+):(.+)$`)
	correctColorRegex = regexp.MustCompile(`^[A-Fa-f0-9]{1,6}$`)
	formatRegex       = regexp.MustCompile(`\.(jpg|jpeg|JPG|JPEG|gif|GIF|png|PNG)`)
)

// Custom responses -----------------------------------------------------------
// Custom response for image
func (r ImageResponse) Apply(req *revel.Request, resp *revel.Response) {

	// FIX:
	// If settings loaded out of actions then revel throws nil pointer, so we
	// load here the first time only
	if font == nil {
		fontPath, _ := revel.Config.String("gummyimage.fontpath")
		font, _ = gummyimage.LoadFont(fontPath)
	}

	resp.WriteHeader(http.StatusOK, "image/png")

	g, _ := gummyimage.NewDefaultGummy(r.sizeX, r.sizeY, r.bgColor)
	g.Font = font

	// Custom text?
	if len(r.text) == 0 {
		g.DrawTextSize(r.fgColor)
	} else {
		g.DrawTextCenter(r.text, r.fgColor)
	}

	b := new(bytes.Buffer)
	g.Get(r.format, b)
	resp.Out.Write(b.Bytes())
}

// Actions --------------------------------------------------------------------
func (c Application) Index() revel.Result {
	return c.Render()
}

func (c Application) CreateImage() revel.Result {

	// Get params by dict because we use this action for 3 different url routes
	// with different url params
	var bgColor, fgColor string
	format, _ := revel.Config.String("gummyimage.format.default")
	text := c.Params.Get("text")

	tmpValues := []string{
		c.Params.Get("size"),
		c.Params.Get("bgcolor"),
		c.Params.Get("fgcolor"),
	}

	// Get format
	for k, i := range tmpValues {
		if f := formatRegex.FindStringSubmatch(i); len(f) > 0 {
			format = f[1]
			tmpValues[k] = formatRegex.ReplaceAllString(i, "")
		}
	}

	x, y, err := getSize(tmpValues[0])
	bgColor, err = colorOk(tmpValues[1])

	if len(tmpValues[2]) > 0 {
		fgColor, err = colorOk(tmpValues[2])
	}

	if err != nil {
		return c.RenderText("Wrong size format")
	}

	// Check limits, don't allow gigantic images :P
	maxY, _ := revel.Config.String("gummyimage.max.height")
	maxX, _ := revel.Config.String("gummyimage.max.width")
	tmx, _ := strconv.Atoi(maxX)
	tmy, _ := strconv.Atoi(maxY)
	if x > tmx || y > tmy {
		return c.RenderText("wow, very big, too image,// Color in HEX format: FAFAFA much pixels")
	}

	return ImageResponse(ImageResponse{x, y, bgColor, fgColor, text, format})
}

// Helpers--------------------------------------------------------------------

// Gets the correct size based on the patern
// Supports:
//  - Predefined sizes (in app.conf)
//  - Aspect sizes: nnnXnn:nn & nn:nnXnnn
//  - Square: nnn
//  - Regular: nnnXnnn & nnnxnnn
func getSize(size string) (x, y int, err error) {

	// Check if is a standard size
	if s, found := revel.Config.String(fmt.Sprintf("size.%v", size)); found {
		size = s
	}

	// Normal size (nnnxnnn, nnnXnnn)
	sizes := regularSizeRegex.FindStringSubmatch(size)
	if len(sizes) > 0 {
		// Check if aspect (nn:nn)

		left := aspectSizeRegex.FindStringSubmatch(sizes[1])
		right := aspectSizeRegex.FindStringSubmatch(sizes[2])

		// If both scale then error
		if len(left) > 0 && len(right) > 0 {
			err = errors.New("Not correct size")
			return

		} else if len(left) > 0 { // nn:nnXnnn
			y, _ = strconv.Atoi(sizes[2])
			tll, _ := strconv.Atoi(left[1])
			tlr, _ := strconv.Atoi(left[2])
			x = y * tll / tlr
		} else if len(right) > 0 { // nnnXnn:nn
			x, _ = strconv.Atoi(sizes[1])
			trl, _ := strconv.Atoi(right[1])
			trr, _ := strconv.Atoi(right[2])
			y = x * trr / trl
		} else { // nnnXnnn
			x, _ = strconv.Atoi(sizes[1])
			y, _ = strconv.Atoi(sizes[2])
		}

	} else { // Square (nnn)
		x, _ = strconv.Atoi(size)
		y = x
	}

	if x == 0 || y == 0 {
		err = errors.New("Not correct size")
	}
	return
}

func colorOk(color string) (newColor string, err error) {

	// Set defaults
	if color == "" {
		newColor, _ = revel.Config.String("gummyimage.bgcolor.default")
		return
	} else if !correctColorRegex.MatchString(color) {
		newColor, _ = revel.Config.String("gummyimage.bgcolor.default")
		err = errors.New("Wrong color format")
		return
	} else {
		switch len(color) {
		case 1:
			newColor = ""
			for i := 0; i < 6; i++ {
				newColor += color
			}
			return
		case 2:
			newColor = fmt.Sprintf("%s%s%s", color, color, color)
			return
		case 3:
			c1 := string(color[0])
			c2 := string(color[1])
			c3 := string(color[2])
			newColor = fmt.Sprintf("%s%s%s%s%s%s", c1, c1, c2, c2, c3, c3)
			return
		}
	}
	newColor = color
	return
}

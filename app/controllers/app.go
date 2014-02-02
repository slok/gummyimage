package controllers

import (
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
	correctColorRegex = regexp.MustCompile(`^[A-Fa-f0-9]{2,6}$`)
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
	g.DrawTextSize(r.fgColor)

	b, _ := g.GetPng()
	resp.Out.Write(b)
}

// Actions --------------------------------------------------------------------
func (c Application) Index() revel.Result {
	return c.Render()
}

func (c Application) CreateImage() revel.Result {

	// Get params by dict because we use this action for 3 different url routes
	// with different url params
	size := c.Params.Get("size")
	bgColor := c.Params.Get("bgcolor")
	fgColor := c.Params.Get("fgcolor")

	bgColor, err := colorOk(bgColor)
	x, y, err := getSize(size)

	if err != nil {
		return c.RenderText("Wrong size format")
	}

	// Check limits, don't allow gigantic images :P
	maxY, _ := revel.Config.String("gummyimage.max.height")
	maxX, _ := revel.Config.String("gummyimage.max.width")
	tmx, _ := strconv.ParseInt(maxX, 10, 0)
	tmy, _ := strconv.ParseInt(maxY, 10, 0)
	if x > int(tmx) || y > int(tmy) {
		return c.RenderText("wow, very big, too image,// Color in HEX format: FAFAFA much pixels")
	}

	return ImageResponse(ImageResponse{x, y, bgColor, fgColor, "", "png"})
}

// Helpers--------------------------------------------------------------------

// Gets the correct size based on the patern
// Supports:
//  - Predefined sizes (in app.conf)
//  - Aspect sizes: nnnXnn:nn & nn:nnXnnn
//  - Square: nnn
//  - Regular: nnnXnnn & nnnxnnn
func getSize(size string) (x, y int, err error) {
	var tx, ty int64

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
			ty, _ = strconv.ParseInt(sizes[2], 10, 0)
			tll, _ := strconv.ParseInt(left[1], 10, 0)
			tlr, _ := strconv.ParseInt(left[2], 10, 0)
			tx = ty * tll / tlr
		} else if len(right) > 0 { // nnnXnn:nn
			tx, _ = strconv.ParseInt(sizes[1], 10, 0)
			trl, _ := strconv.ParseInt(right[1], 10, 0)
			trr, _ := strconv.ParseInt(right[2], 10, 0)
			ty = tx * trr / trl
		} else { // nnnXnnn
			tx, _ = strconv.ParseInt(sizes[1], 10, 0)
			ty, _ = strconv.ParseInt(sizes[2], 10, 0)
		}

		x = int(tx)
		y = int(ty)

	} else { // Square (nnn)
		tx, _ := strconv.ParseInt(size, 10, 0)
		x = int(tx)
		y = x
	}

	if x == 0 || y == 0 {
		err = errors.New("Not correct size")
	}
	return
}

func colorOk(color string) (bgColor string, err error) {

	// Set defaults
	if color == "" {
		bgColor, _ = revel.Config.String("gummyimage.bgcolor.default")
		return
	} else if !correctColorRegex.MatchString(color) {
		bgColor, _ = revel.Config.String("gummyimage.bgcolor.default")
		err = errors.New("Wrong color format")
		return
	} else {
		switch len(color) {
		case 2:
			bgColor = fmt.Sprintf("%s%s%s", color, color, color)
			return
		case 3:
			c1 := string(color[0])
			c2 := string(color[1])
			c3 := string(color[2])
			bgColor = fmt.Sprintf("%s%s%s%s%s%s", c1, c1, c2, c2, c3, c3)
			return
		}
	}
	bgColor = color
	return
}

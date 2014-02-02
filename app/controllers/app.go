package controllers

import (
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
	font             *truetype.Font
	regularSizeRegex = regexp.MustCompile(`^(.+)[xX](.+)$`)
	aspectSizeRegex  = regexp.MustCompile(`^(.+):(.+)$`)
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

	//TODO: Random color for now
	g, _ := gummyimage.NewDefaultGummy(r.sizeX, r.sizeY)
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

	var x, y int
	var tx, ty int64

	// Get correct size flow

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
			//return error
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

	// Check limits, don't allow gigantic images :P

	return ImageResponse(ImageResponse{x, y, bgColor, fgColor, "", "png"})
}

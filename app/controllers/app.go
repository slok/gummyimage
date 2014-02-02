package controllers

import (
	"net/http"

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
	font *truetype.Font
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
	//size := c.Params.Get("size")
	//bgColor := c.Params.Get("bgcolor")
	//fgColor := c.Params.Get("fgcolor")

	return ImageResponse(ImageResponse{600, 600, "", "", "", "png"})
}

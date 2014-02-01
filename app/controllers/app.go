package controllers

import (
	"fmt"

	"github.com/robfig/revel"
)

type Application struct {
	*revel.Controller
}

func (c Application) Index() revel.Result {
	return c.Render()
}

func (c Application) CreateImage() revel.Result {

	// Get params by dict because we use this action for 3 different url routes
	// with different url params
	size := c.Params.Get("size")
	bgColor := c.Params.Get("bgcolor")
	fgColor := c.Params.Get("fgcolor")

	fmt.Println(size)
	fmt.Println(bgColor)
	fmt.Println(fgColor)

	return c.Render()
}

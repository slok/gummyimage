package controllers

import (
	"github.com/robfig/revel"
)

type Application struct {
}

func (c Application) Index() revel.Result {
	return c.Render()
}

func (c Application) CreateImage(size, bgcolor, fgcolor string) revel.Result {
	return c.Render()
}

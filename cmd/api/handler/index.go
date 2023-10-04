package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type IndexHandler interface {
	Index(echo.Context) error
}

type index struct{}

func NewIndex() IndexHandler {
	return &index{}
}

func (i *index) Index(c echo.Context) (err error) {
	return c.HTML(http.StatusOK, "<h1 style='text-align:center;'>Welcome, TODO API!</h1>")
}

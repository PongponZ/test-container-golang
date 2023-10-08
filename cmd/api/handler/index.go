package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) (err error) {
	return c.HTML(http.StatusOK, "<h1 style='text-align:center;'>Welcome, TODO API!</h1>")
}

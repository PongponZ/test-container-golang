package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type TodoHandler interface {
	Create(echo.Context) error
}

type todo struct{}

func NewTodo() TodoHandler {
	return &todo{}
}

func (t *todo) Create(c echo.Context) (err error) {
	return c.String(http.StatusOK, "Hello, World!")
}

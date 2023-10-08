package handler

import (
	"net/http"

	"github.com/PongponZ/test-container-golang/pkg/entity"
	"github.com/PongponZ/test-container-golang/pkg/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoHandler interface {
	Create(echo.Context) error
	Gets(echo.Context) error
	Update(echo.Context) error
	Delete(echo.Context) error
}

type todo struct {
	usecase usecase.TodoUsecase
}

func NewTodo(usecase usecase.TodoUsecase) TodoHandler {
	return &todo{
		usecase: usecase,
	}
}

func (t *todo) Create(c echo.Context) (err error) {
	var payload TodoPayload
	if err = c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	todo := entity.Todo{
		Title:       payload.Title,
		Description: payload.Description,
	}

	insertID, err := t.usecase.Create(todo)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.String(http.StatusOK, insertID)
}

func (t *todo) Gets(c echo.Context) (err error) {
	todos, err := t.usecase.Gets()
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusOK, todos)
}

func (t *todo) Update(c echo.Context) (err error) {
	var payload TodoUpdatePayload
	if err = c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	id, err := primitive.ObjectIDFromHex(payload.ID)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	todo := entity.Todo{
		ID:          id,
		Title:       payload.Title,
		Description: payload.Description,
	}

	if err = t.usecase.Update(todo); err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.String(http.StatusOK, "Update Success")
}

func (t *todo) Delete(c echo.Context) (err error) {
	id := c.Param("id")
	if id == "" {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	if err = t.usecase.Delete(id); err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.String(http.StatusOK, "Delete Success")
}

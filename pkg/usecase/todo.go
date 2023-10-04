package usecase

import "github.com/PongponZ/test-container-golang/pkg/entity"

type TodoUsecase interface {
	Create(entity.Todo) error
}

type todo struct{}

func NewToDo() TodoUsecase {
	return &todo{}
}

func (t *todo) Create(task entity.Todo) (err error) {
	return nil
}

package usecase

import (
	"github.com/PongponZ/test-container-golang/pkg/entity"
	"github.com/PongponZ/test-container-golang/pkg/repository"
)

type TodoUsecase interface {
	Create(entity.Todo) (string, error)
	Gets() ([]entity.Todo, error)
	Update(entity.Todo) error
	Delete(string) error
}

type todo struct {
	repo repository.TodoRepository
}

func NewToDo(repo repository.TodoRepository) TodoUsecase {
	return &todo{
		repo: repo,
	}
}

func (t *todo) Create(task entity.Todo) (string, error) {
	return t.repo.Create(task)
}

func (t *todo) Gets() ([]entity.Todo, error) {
	return t.repo.Gets()
}

func (t *todo) Update(task entity.Todo) error {
	return t.repo.Update(task)
}

func (t *todo) Delete(id string) error {
	return t.repo.Delete(id)
}

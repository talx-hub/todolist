package service

import (
	"github.com/talx-hub/todolist/internal/model"
	"github.com/talx-hub/todolist/internal/repo"
)

type TodoList struct {
	repo repo.Repository
}

func New(repo repo.Repository) TodoList {
	return TodoList{repo: repo}
}

func (l *TodoList) GetAll() []model.Task {
	return l.repo.GetAll()
}

func (l *TodoList) Get(id string) (model.Task, error) {
	return l.repo.Get(id)
}

func (l *TodoList) Post(task model.Task) error {
	return l.repo.Post(task)
}

func (l *TodoList) Delete(id string) error {
	return l.repo.Delete(id)
}

package repo

import (
	"fmt"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/model"
)

type Repository interface {
	GetAll() []model.Task
	Get(id string) (model.Task, error)
	Put(task model.Task)
	Delete(id string) error
}

type InMemoryTasks map[string]model.Task

func (t InMemoryTasks) GetAll() []model.Task {
	tasks := make([]model.Task, 0)
	for _, task := range t {
		tasks = append(tasks, task)
	}
	return tasks
}

func (t InMemoryTasks) Get(id string) (model.Task, error) {
	task, ok := t[id]
	if !ok {
		return model.Task{}, fmt.Errorf("id not found: %s", id)
	}
	return task, nil
}

func (t InMemoryTasks) Put(task model.Task) {
	id := task.ID
	t[id] = task
}

func (t InMemoryTasks) Delete(id string) error {
	if _, ok := t[id]; !ok {
		return fmt.Errorf("id not found %s", id)
	}
	delete(t, id)
	return nil
}

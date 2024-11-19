package main

import (
	"fmt"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/api"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/repo"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func prepare() *api.Handler {
	tasks := repo.InMemoryTasks{
		"1": {
			ID:          "1",
			Description: "Сделать финальное задание темы REST API",
			Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
			Applications: []string{
				"VS Code",
				"Terminal",
				"git",
			},
		},
		"2": {
			ID:          "2",
			Description: "Протестировать финальное задание с помощью Postman",
			Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
			Applications: []string{
				"VS Code",
				"Terminal",
				"git",
				"Postman",
			},
		},
	}

	todoList := service.New(tasks)
	return api.New(todoList)
}

func TaskRouter() chi.Router {
	handlers := prepare()
	router := chi.NewRouter()
	router.Get("/tasks", handlers.GetTasks)
	router.Get("/tasks/{id}", handlers.GetTask)
	router.Post("/tasks", handlers.PostTask)
	router.Delete("/tasks/{id}", handlers.DeleteTask)
	return router
}

func main() {
	if err := http.ListenAndServe(":8080", TaskRouter()); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}

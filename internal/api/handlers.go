package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Yandex-Practicum/go-rest-api-homework/internal/model"
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/service"
)

type Handler struct {
	service service.TodoList
}

func New(service service.TodoList) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(h.service.GetAll()); err != nil {
		http.Error(
			w,
			fmt.Sprintf("cannot encode tasks: %s", err.Error()),
			http.StatusInternalServerError)
	}
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, err := h.service.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(task); err != nil {
		http.Error(
			w,
			fmt.Sprintf("cannot encode task: %s", err.Error()),
			http.StatusInternalServerError)
		// да, по ТЗ просят BadRequest, но кажется правильным вернуть ошибку сервера
	}
}

func (h *Handler) PostTask(w http.ResponseWriter, r *http.Request) {
	task := model.Task{}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(
			w,
			fmt.Sprintf("cannot read body: %s", err.Error()),
			http.StatusBadRequest)
		return
	}

	if err := h.service.Post(task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

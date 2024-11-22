package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/talx-hub/todolist/internal/model"
	"github.com/talx-hub/todolist/internal/repo"
	"github.com/talx-hub/todolist/internal/service"
)

func prepareHandlers() []*Handler {
	tasksEmpty := repo.InMemoryTasks{}
	todoListEmpty := service.New(tasksEmpty)
	handlerEmpty := New(todoListEmpty)

	tasksSingle := repo.InMemoryTasks{
		"1": {
			ID:           "1",
			Description:  "",
			Note:         "",
			Applications: []string{},
		},
	}

	todoListSingle := service.New(tasksSingle)
	handlerSingle := New(todoListSingle)

	tasksMany := repo.InMemoryTasks{
		"1": {
			ID:           "1",
			Description:  "",
			Note:         "",
			Applications: []string{},
		},
		"2": {
			ID:           "2",
			Description:  "",
			Note:         "",
			Applications: []string{},
		},
	}

	todoListMany := service.New(tasksMany)
	handlerMany := New(todoListMany)
	return []*Handler{handlerEmpty, handlerSingle, handlerMany}
}

func testRequest(t *testing.T,
	method string, handler http.HandlerFunc,
	endpoint, id string, body io.Reader) (*http.Response, string) {

	r := httptest.NewRequest(method, endpoint, body)
	w := httptest.NewRecorder()
	chiCtx := chi.NewRouteContext()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
	chiCtx.URLParams.Add("id", id)
	handler(w, r)

	response := w.Result()
	defer func() {
		err := response.Body.Close()
		require.NoError(t, err)
	}()
	resBody, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return response, string(resBody)
}

func TestGetTasks(t *testing.T) {
	type want struct {
		status   int
		response string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "get from empty",
			want: want{
				status:   http.StatusOK,
				response: `[]`,
			},
		},
		{
			name: "get single",
			want: want{
				status: http.StatusOK,
				response: `
[{"id": "1", "description": "", "note": "", "applications": []}]`},
		},
		{
			name: "get many",
			want: want{
				status: http.StatusOK,
				response: `[
{"id": "1", "description": "", "note": "", "applications": []},
{"id": "2", "description": "", "note": "", "applications": []}]`,
			},
		},
	}

	handlers := prepareHandlers()
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, got := testRequest(
				t, http.MethodGet, handlers[i].GetTasks,
				"/tasks", "", http.NoBody)
			assert.Equal(t, test.want.status, resp.StatusCode)
			tasksWant := make([]model.Task, 0)
			err := json.Unmarshal([]byte(test.want.response), &tasksWant)
			require.NoError(t, err)
			tasksGot := make([]model.Task, 0)
			err = json.Unmarshal([]byte(got), &tasksGot)
			require.NoError(t, err)
			assert.ElementsMatch(t, tasksWant, tasksGot)
		})
	}
}

func TestGetTask(t *testing.T) {
	type want struct {
		status   int
		response string
	}
	tests := []struct {
		name     string
		endpoint string
		id       string
		want     want
	}{
		{
			name:     "positive #1",
			endpoint: "/tasks",
			id:       "1",
			want: want{
				status: http.StatusOK,
				response: ` 
{"id": "1", "description": "", "note": "", "applications": []}`,
			},
		},
		{
			name:     "negative #1",
			endpoint: "/tasks/5",
			id:       "5",
			want: want{
				status:   http.StatusBadRequest,
				response: "id not found: 5\n",
			},
		},
	}

	handlers := prepareHandlers()
	const IDTasksMany = 2
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, got := testRequest(
				t, http.MethodGet, handlers[IDTasksMany].GetTask,
				test.endpoint, test.id, http.NoBody)
			assert.Equal(t, test.want.status, resp.StatusCode)
			if resp.StatusCode == http.StatusOK {
				assert.JSONEq(t, test.want.response, got)
			} else {
				assert.Equal(t, test.want.response, got)
			}
		})
	}
}

func TestPostTask(t *testing.T) {
	type want struct {
		status   int
		response string
	}
	tests := []struct {
		name     string
		endpoint string
		body     string
		want     want
	}{
		{
			name:     "positive #1",
			endpoint: "/tasks",
			body: `
{"id": "1", "description": "", "note": "", "applications": []}`,
			want: want{
				status:   http.StatusCreated,
				response: "",
			},
		},
		{
			name:     "negative #1",
			endpoint: "/tasks",
			body: `
{"id": "1", "description": "", "note": "", "applications": []}`,
			want: want{
				status:   http.StatusBadRequest,
				response: "task already exists\n",
			},
		},
	}

	handlers := prepareHandlers()
	const IDTasksEmpty = 0
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, got := testRequest(
				t, http.MethodPost, handlers[IDTasksEmpty].PostTask,
				test.endpoint, "", strings.NewReader(test.body))
			assert.Equal(t, test.want.status, resp.StatusCode)
			assert.Equal(t, test.want.response, got)
		})
	}
}

func TestDeleteTask(t *testing.T) {
	type want struct {
		status   int
		response string
	}
	tests := []struct {
		name     string
		endpoint string
		id       string
		want     want
	}{
		{
			name:     "negative #1",
			endpoint: "/tasks",
			id:       "3",
			want: want{
				status:   http.StatusBadRequest,
				response: "id not found: 3\n",
			},
		},
		{
			name:     "positive #1",
			endpoint: "/tasks",
			id:       "1",
			want: want{
				status:   http.StatusOK,
				response: "",
			},
		},
		{
			name:     "positive #2",
			endpoint: "/tasks",
			id:       "2",
			want: want{
				status:   http.StatusOK,
				response: "",
			},
		},
		{
			name:     "negative #2",
			endpoint: "/tasks",
			id:       "1",
			want: want{
				status:   http.StatusBadRequest,
				response: "id not found: 1\n",
			},
		},
		{
			name:     "negative #3",
			endpoint: "/tasks",
			id:       "2",
			want: want{
				status:   http.StatusBadRequest,
				response: "id not found: 2\n",
			},
		},
	}

	handlers := prepareHandlers()
	const IDTasksMany = 2
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, got := testRequest(
				t, http.MethodDelete, handlers[IDTasksMany].DeleteTask,
				test.endpoint, test.id, http.NoBody)
			assert.Equal(t, test.want.status, resp.StatusCode)
			assert.Equal(t, test.want.response, got)
		})
	}
}

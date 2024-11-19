package repo

import (
	"github.com/Yandex-Practicum/go-rest-api-homework/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestGetAll(t *testing.T) {
	tests := []struct {
		name   string
		length int
		want   []model.Task
	}{
		{
			name:   "get all from empty repo",
			length: 0,
			want:   []model.Task{},
		},
		{
			name:   "get one",
			length: 1,
			want: []model.Task{
				{
					ID:           "1",
					Description:  "descr1",
					Note:         "note1",
					Applications: []string{"11", "12", "13"},
				}},
		},
		{
			name:   "get many",
			length: 3,
			want: []model.Task{
				{
					ID:           "1",
					Description:  "descr1",
					Note:         "note1",
					Applications: []string{"11", "12", "13"},
				},
				{
					ID:           "2",
					Description:  "descr2",
					Note:         "note2",
					Applications: []string{"21", "22", "23"},
				},
				{
					ID:           "3",
					Description:  "descr3",
					Note:         "note3",
					Applications: []string{"31", "32", "33"},
				},
			},
		},
	}

	repoEmpty := InMemoryTasks{}
	repoLonely := InMemoryTasks{
		"1": {
			ID:           "1",
			Description:  "descr1",
			Note:         "note1",
			Applications: []string{"11", "12", "13"},
		},
	}
	repoMany := InMemoryTasks{
		"1": {
			ID:           "1",
			Description:  "descr1",
			Note:         "note1",
			Applications: []string{"11", "12", "13"},
		},
		"2": {
			ID:           "2",
			Description:  "descr2",
			Note:         "note2",
			Applications: []string{"21", "22", "23"},
		},
		"3": {
			ID:           "3",
			Description:  "descr3",
			Note:         "note3",
			Applications: []string{"31", "32", "33"},
		},
	}
	repos := []*InMemoryTasks{&repoEmpty, &repoLonely, &repoMany}
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.NotNil(t, repos[i])
			tasks := repos[i].GetAll()
			assert.Equal(t, test.length, len(tasks))
			assert.ElementsMatch(t, test.want, tasks)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		want    model.Task
		errWant bool
	}{
		{
			name: "positive test: get existing",
			id:   "1",
			want: model.Task{
				ID:           "1",
				Description:  "descr1",
				Note:         "note1",
				Applications: []string{"11", "12", "13"},
			},
			errWant: false,
		},
		{
			name:    "negative test #1, get non existing",
			id:      "4",
			want:    model.Task{},
			errWant: true,
		},
		{
			name:    "negative test #2: get from empty",
			id:      "1",
			want:    model.Task{},
			errWant: true,
		},
	}

	repoEmpty := InMemoryTasks{}
	repoMany := InMemoryTasks{
		"1": {
			ID:           "1",
			Description:  "descr1",
			Note:         "note1",
			Applications: []string{"11", "12", "13"},
		},
		"2": {
			ID:           "2",
			Description:  "descr2",
			Note:         "note2",
			Applications: []string{"21", "22", "23"},
		},
		"3": {
			ID:           "3",
			Description:  "descr3",
			Note:         "note3",
			Applications: []string{"31", "32", "33"},
		},
	}
	repos := []*InMemoryTasks{&repoMany, &repoMany, &repoEmpty}
	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.NotNil(t, repos[i])
			task, err := repos[i].Get(test.id)
			if test.errWant {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, task)
		})
	}
}

func TestPost(t *testing.T) {
	tests := []struct {
		name       string
		task       model.Task
		lengthWant int
		want       InMemoryTasks
		errWant    bool
	}{
		{
			name: "positive test #1: post new to empty",
			task: model.Task{
				ID:           "1",
				Description:  "descr1",
				Note:         "note1",
				Applications: []string{"11", "12", "13"},
			},
			lengthWant: 1,
			want: InMemoryTasks{
				"1": {
					ID:           "1",
					Description:  "descr1",
					Note:         "note1",
					Applications: []string{"11", "12", "13"},
				},
			},
			errWant: false,
		},
		{
			name: "positive test #2: post new",
			task: model.Task{
				ID:           "2",
				Description:  "descr2",
				Note:         "note2",
				Applications: []string{"21", "22", "23"},
			},
			lengthWant: 2,
			want: InMemoryTasks{
				"1": {
					ID:           "1",
					Description:  "descr1",
					Note:         "note1",
					Applications: []string{"11", "12", "13"},
				},
				"2": {
					ID:           "2",
					Description:  "descr2",
					Note:         "note2",
					Applications: []string{"21", "22", "23"},
				},
			},
			errWant: false,
		},
		{
			name: "negative test #1: post existing",
			task: model.Task{
				ID:           "2",
				Description:  "descr2",
				Note:         "note2",
				Applications: []string{"21", "22", "23"},
			},
			lengthWant: 2,
			want: InMemoryTasks{
				"1": {
					ID:           "1",
					Description:  "descr1",
					Note:         "note1",
					Applications: []string{"11", "12", "13"},
				},
				"2": {
					ID:           "2",
					Description:  "descr2",
					Note:         "note2",
					Applications: []string{"21", "22", "23"},
				},
			},
			errWant: true,
		},
	}

	repo := InMemoryTasks{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := repo.Post(test.task)
			if test.errWant {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.lengthWant, len(repo))
			assert.True(t, reflect.DeepEqual(test.want, repo))
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name       string
		ID         string
		lengthWant int
		want       InMemoryTasks
		errWant    bool
	}{
		{
			name:       "positive test #1: delete existing",
			ID:         "1",
			lengthWant: 1,
			want: InMemoryTasks{
				"2": {
					ID:           "2",
					Description:  "descr2",
					Note:         "note2",
					Applications: []string{"21", "22", "23"},
				},
			},
			errWant: false,
		},
		{
			name:       "negative test #1: delete the same",
			ID:         "1",
			lengthWant: 1,
			want: InMemoryTasks{
				"2": {
					ID:           "2",
					Description:  "descr2",
					Note:         "note2",
					Applications: []string{"21", "22", "23"},
				},
			},
			errWant: true,
		},
		{
			name:       "negative test #2: delete non existing",
			ID:         "3",
			lengthWant: 1,
			want: InMemoryTasks{
				"2": {
					ID:           "2",
					Description:  "descr2",
					Note:         "note2",
					Applications: []string{"21", "22", "23"},
				},
			},
			errWant: true,
		},
		{
			name:       "positive test #2: delete existing",
			ID:         "2",
			lengthWant: 0,
			want:       InMemoryTasks{},
			errWant:    false,
		},
		{
			name:       "negative test #3: delete from empty",
			ID:         "1",
			lengthWant: 0,
			want:       InMemoryTasks{},
			errWant:    true,
		},
	}

	repo := InMemoryTasks{
		"1": {
			ID:           "1",
			Description:  "descr1",
			Note:         "note1",
			Applications: []string{"11", "12", "13"},
		},
		"2": {
			ID:           "2",
			Description:  "descr2",
			Note:         "note2",
			Applications: []string{"21", "22", "23"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := repo.Delete(test.ID)
			if test.errWant {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.lengthWant, len(repo))
			assert.True(t, reflect.DeepEqual(test.want, repo))
		})
	}
}

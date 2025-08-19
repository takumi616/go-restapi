package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/takumi616/go-restapi/domain"
)

func TestInsert(t *testing.T) {
	type expected struct {
		task *domain.Task
		err  error
	}

	testTable := map[string]struct {
		input     *domain.Task
		mockSetup func(sqlmock.Sqlmock)
		expected  expected
	}{
		"ok": {
			input: &domain.Task{
				Title:       "Test Title",
				Description: "Test Description",
				Status:      false,
			},
			mockSetup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "status"}).
					AddRow("6a30b9b0-18bf-47b4-bd23-d72726864def", "Test Title", "Test Description", false)

				m.ExpectQuery(regexp.QuoteMeta(
					"INSERT INTO tasks(title, description, status) VALUES($1, $2, $3) RETURNING *",
				)).
					WithArgs("Test Title", "Test Description", false).
					WillReturnRows(rows)
			},
			expected: expected{
				task: &domain.Task{
					Id:          "6a30b9b0-18bf-47b4-bd23-d72726864def",
					Title:       "Test Title",
					Description: "Test Description",
					Status:      false,
				},
				err: nil,
			},
		},
		"duplicateError": {
			input: &domain.Task{
				Title:       "Duplicate Title",
				Description: "Test Description",
				Status:      false,
			},
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(regexp.QuoteMeta(
					"INSERT INTO tasks(title, description, status) VALUES($1, $2, $3) RETURNING *",
				)).
					WithArgs("Duplicate Title", "Test Description", false).
					WillReturnError(
						errors.New(
							"pq: duplicate key value violates unique constraint \"tasks_title_key\"",
						),
					)
			},
			expected: expected{
				task: nil,
				err: errors.New(
					"failed to execute query: pq: duplicate key value violates unique constraint \"tasks_title_key\"",
				),
			},
		},
	}

	for n, tt := range testTable {
		tt := tt
		t.Run(n, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			repo := &TaskRepository{Db: db}
			result, err := repo.Insert(context.Background(), tt.input)

			if n == "ok" {
				assert.Equal(t, tt.expected.task, result)
				assert.Nil(t, err)
			} else if n == "duplicateError" {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expected.err.Error())
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSelectAll(t *testing.T) {
	type expected struct {
		taskList []*domain.Task
	}

	testTable := map[string]struct {
		mockSetup func(sqlmock.Sqlmock)
		expected  expected
	}{
		"ok": {
			mockSetup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "status"}).
					AddRow("6a30b9b0-18bf-47b4-bd23-d72726864def", "Test Title", "Test Description", false).
					AddRow("3e440171-0921-4c88-a7ec-13f4cdab0d69", "Test Title2", "Test Description2", false)

				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM tasks",
				)).WillReturnRows(rows)
			},
			expected: expected{
				taskList: []*domain.Task{
					{
						Id:          "6a30b9b0-18bf-47b4-bd23-d72726864def",
						Title:       "Test Title",
						Description: "Test Description",
						Status:      false,
					},
					{
						Id:          "3e440171-0921-4c88-a7ec-13f4cdab0d69",
						Title:       "Test Title2",
						Description: "Test Description2",
						Status:      false,
					},
				},
			},
		},
		"empty": {
			mockSetup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "status"})

				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM tasks",
				)).WillReturnRows(rows)
			},
			expected: expected{
				taskList: nil,
			},
		},
	}

	for n, tt := range testTable {
		tt := tt
		t.Run(n, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock)

			repo := &TaskRepository{Db: db}
			result, err := repo.SelectAll(context.Background())

			assert.Equal(t, tt.expected.taskList, result)
			assert.Nil(t, err)

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSelectById(t *testing.T) {
	type expected struct {
		task *domain.Task
		err  error
	}

	testTable := map[string]struct {
		id        string
		mockSetup func(sqlmock.Sqlmock, string)
		expected  expected
	}{
		"ok": {
			id: "6a30b9b0-18bf-47b4-bd23-d72726864def",
			mockSetup: func(m sqlmock.Sqlmock, id string) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "status"}).
					AddRow("6a30b9b0-18bf-47b4-bd23-d72726864def", "Test Title", "Test Description", false)

				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM tasks WHERE id = $1",
				)).WithArgs(id).WillReturnRows(rows)
			},
			expected: expected{
				task: &domain.Task{
					Id:          "6a30b9b0-18bf-47b4-bd23-d72726864def",
					Title:       "Test Title",
					Description: "Test Description",
					Status:      false,
				},
				err: nil,
			},
		},
		"invalidId": {
			id: "abc123",
			mockSetup: func(m sqlmock.Sqlmock, id string) {
				sqlmock.NewRows([]string{"id", "title", "description", "status"}).
					AddRow("6a30b9b0-18bf-47b4-bd23-d72726864def", "Test Title", "Test Description", false)

				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM tasks WHERE id = $1",
				)).WithArgs(id).WillReturnError(errors.New("pq: invalid input syntax for type uuid: \"abc123\""))
			},
			expected: expected{
				task: nil,
				err:  errors.New("failed to copy columns: pq: invalid input syntax for type uuid: \"abc123\""),
			},
		},
	}

	for n, tt := range testTable {
		tt := tt
		t.Run(n, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock, tt.id)

			repo := &TaskRepository{Db: db}
			result, err := repo.SelectById(context.Background(), tt.id)

			if n == "ok" {
				assert.Equal(t, tt.expected.task, result)
				assert.Nil(t, err)
			} else if n == "invalidId" {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expected.err.Error())
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdate(t *testing.T) {
	type expected struct {
		task *domain.Task
		err  error
	}

	testTable := map[string]struct {
		id        string
		input     *domain.Task
		mockSetup func(sqlmock.Sqlmock, string, *domain.Task)
		expected  expected
	}{
		"ok": {
			id: "6a30b9b0-18bf-47b4-bd23-d72726864def",
			input: &domain.Task{
				Description: "Update Test Description",
				Status:      true,
			},
			mockSetup: func(m sqlmock.Sqlmock, id string, input *domain.Task) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "status"}).
					AddRow("6a30b9b0-18bf-47b4-bd23-d72726864def", "Test Title", "Update Test Description", true)

				m.ExpectQuery(regexp.QuoteMeta(
					"UPDATE tasks SET description=$1, status=$2 WHERE id=$3 RETURNING *",
				)).
					WithArgs(input.Description, input.Status, id).
					WillReturnRows(rows)
			},
			expected: expected{
				task: &domain.Task{
					Id:          "6a30b9b0-18bf-47b4-bd23-d72726864def",
					Title:       "Test Title",
					Description: "Update Test Description",
					Status:      true,
				},
				err: nil,
			},
		},
		"invalidId": {
			id: "abc123",
			input: &domain.Task{
				Description: "Update Test Description",
				Status:      true,
			},
			mockSetup: func(m sqlmock.Sqlmock, id string, input *domain.Task) {
				sqlmock.NewRows([]string{"id", "title", "description", "status"}).
					AddRow("6a30b9b0-18bf-47b4-bd23-d72726864def", "Test Title", "Test Description", false)

				m.ExpectQuery(regexp.QuoteMeta(
					"UPDATE tasks SET description=$1, status=$2 WHERE id=$3 RETURNING *",
				)).
					WithArgs(input.Description, input.Status, id).
					WillReturnError(errors.New("pq: invalid input syntax for type uuid: \"abc123\""))
			},
			expected: expected{
				task: nil,
				err:  errors.New("failed to copy columns: pq: invalid input syntax for type uuid: \"abc123\""),
			},
		},
	}

	for n, tt := range testTable {
		tt := tt
		t.Run(n, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock, tt.id, tt.input)

			repo := &TaskRepository{Db: db}
			result, err := repo.Update(context.Background(), tt.id, tt.input)

			if n == "ok" {
				assert.Equal(t, tt.expected.task, result)
				assert.Nil(t, err)
			} else if n == "invalidId" {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expected.err.Error())
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDelete(t *testing.T) {
	type expected struct {
		task *domain.Task
		err  error
	}

	testTable := map[string]struct {
		id        string
		mockSetup func(sqlmock.Sqlmock, string)
		expected  expected
	}{
		"ok": {
			id: "6a30b9b0-18bf-47b4-bd23-d72726864def",
			mockSetup: func(m sqlmock.Sqlmock, id string) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow("6a30b9b0-18bf-47b4-bd23-d72726864def")

				m.ExpectQuery(regexp.QuoteMeta(
					"DELETE FROM tasks WHERE id=$1 RETURNING id",
				)).
					WithArgs(id).
					WillReturnRows(rows)
			},
			expected: expected{
				task: &domain.Task{
					Id: "6a30b9b0-18bf-47b4-bd23-d72726864def",
				},
				err: nil,
			},
		},
		"invalidId": {
			id: "abc123",
			mockSetup: func(m sqlmock.Sqlmock, id string) {
				sqlmock.NewRows([]string{"id"}).
					AddRow("6a30b9b0-18bf-47b4-bd23-d72726864def")

				m.ExpectQuery(regexp.QuoteMeta(
					"DELETE FROM tasks WHERE id=$1 RETURNING id",
				)).
					WithArgs(id).
					WillReturnError(errors.New("pq: invalid input syntax for type uuid: \"abc123\""))
			},
			expected: expected{
				task: nil,
				err:  errors.New("failed to copy columns: pq: invalid input syntax for type uuid: \"abc123\""),
			},
		},
	}

	for n, tt := range testTable {
		tt := tt
		t.Run(n, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock, tt.id)

			repo := &TaskRepository{Db: db}
			result, err := repo.Delete(context.Background(), tt.id)

			if n == "ok" {
				assert.Equal(t, tt.expected.task, result)
				assert.Nil(t, err)
			} else if n == "invalidId" {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expected.err.Error())
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

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

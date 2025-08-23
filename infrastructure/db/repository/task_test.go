package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/takumi616/go-restapi/domain"
	"github.com/takumi616/go-restapi/infrastructure/db/repository/model"
	customError "github.com/takumi616/go-restapi/shared/error"
)

func TestInsert(t *testing.T) {
	type expected struct {
		task *domain.Task
		err  error
	}

	testTable := map[string]struct {
		input     *domain.Task
		mockSetup func(sqlmock.Sqlmock, *model.InsertTaskParam)
		expected  expected
	}{
		"Ok": {
			input: &domain.Task{
				Title:       "Test Title",
				Description: "Test Description",
				Status:      false,
			},
			mockSetup: func(m sqlmock.Sqlmock, param *model.InsertTaskParam) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "status"}).
					AddRow("6a30b9b0-18bf-47b4-bd23-d72726864def", param.Title, param.Description, param.Status)

				m.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO tasks(title, description, status)
					VALUES($1, $2, $3)
					RETURNING id, title, description, status`,
				)).
					WithArgs(param.Title, param.Description, param.Status).
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
		"DuplicateError": {
			input: &domain.Task{
				Title:       "Duplicate Title",
				Description: "Test Description",
				Status:      false,
			},
			mockSetup: func(m sqlmock.Sqlmock, param *model.InsertTaskParam) {
				m.ExpectQuery(regexp.QuoteMeta(
					`INSERT INTO tasks(title, description, status)
					VALUES($1, $2, $3)
					RETURNING id, title, description, status`,
				)).
					WithArgs(param.Title, param.Description, param.Status).
					WillReturnError(
						errors.New(
							"pq: duplicate key value violates unique constraint \"tasks_title_key\"",
						),
					)
			},
			expected: expected{
				task: nil,
				err: errors.New(
					customError.ErrInternalServerError.Error(),
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

			tt.mockSetup(mock, model.ToInsertTaskParam(tt.input))

			repo := &TaskRepository{Db: db}
			result, err := repo.Insert(context.Background(), tt.input)

			if tt.expected.err != nil {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expected.err.Error())
			} else {
				assert.Equal(t, tt.expected.task, result)
				assert.Nil(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSelectAll(t *testing.T) {
	type expected struct {
		taskList []*domain.Task
		err      error
	}

	testTable := map[string]struct {
		mockSetup func(sqlmock.Sqlmock)
		expected  expected
	}{
		"Ok": {
			mockSetup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "status"}).
					AddRow("6a30b9b0-18bf-47b4-bd23-d72726864def", "Test Title", "Test Description", false).
					AddRow("3e440171-0921-4c88-a7ec-13f4cdab0d69", "Test Title2", "Test Description2", false)

				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT id, title, description, status FROM tasks",
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
				err: nil,
			},
		},
		"Empty": {
			mockSetup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "status"})

				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT id, title, description, status FROM tasks",
				)).WillReturnRows(rows)
			},
			expected: expected{
				taskList: []*domain.Task{},
				err:      nil,
			},
		},
		"InternalServerErr": {
			mockSetup: func(m sqlmock.Sqlmock) {
				sqlmock.NewRows([]string{"id", "title", "description", "status"})

				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT id, title, description, status FROM tasks",
				)).WillReturnError(errors.New("sql: expected 4 destination arguments in Scan, not 3"))
			},
			expected: expected{
				taskList: nil,
				err:      customError.ErrInternalServerError,
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

			if tt.expected.err != nil {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expected.err.Error())
			} else {
				assert.Equal(t, tt.expected.taskList, result)
				assert.Nil(t, err)
			}

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
		"Ok": {
			id: "6a30b9b0-18bf-47b4-bd23-d72726864def",
			mockSetup: func(m sqlmock.Sqlmock, id string) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "status"}).
					AddRow("6a30b9b0-18bf-47b4-bd23-d72726864def", "Test Title", "Test Description", false)

				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT id, title, description, status FROM tasks WHERE id = $1",
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
		"NotFound": {
			id: "3e440171-0921-4c88-a7ec-13f4cdab0d69",
			mockSetup: func(m sqlmock.Sqlmock, id string) {
				sqlmock.NewRows([]string{"id", "title", "description", "status"}).
					AddRow("6a30b9b0-18bf-47b4-bd23-d72726864def", "Test Title", "Test Description", false)

				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT id, title, description, status FROM tasks WHERE id = $1",
				)).WithArgs(id).WillReturnError(sql.ErrNoRows)
			},
			expected: expected{
				task: nil,
				err:  customError.ErrNotFound,
			},
		},
		"InvalidId": {
			id: "abc123",
			mockSetup: func(m sqlmock.Sqlmock, id string) {
				sqlmock.NewRows([]string{"id", "title", "description", "status"}).
					AddRow("6a30b9b0-18bf-47b4-bd23-d72726864def", "Test Title", "Test Description", false)

				m.ExpectQuery(regexp.QuoteMeta(
					"SELECT id, title, description, status FROM tasks WHERE id = $1",
				)).WithArgs(id).WillReturnError(errors.New("pq: invalid input syntax for type uuid: \"abc123\""))
			},
			expected: expected{
				task: nil,
				err:  customError.ErrInternalServerError,
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

			if tt.expected.err != nil {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expected.err.Error())
			} else {
				assert.Equal(t, tt.expected.task, result)
				assert.Nil(t, err)
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
		mockSetup func(sqlmock.Sqlmock, string, *model.UpdateTaskParam)
		expected  expected
	}{
		"Ok": {
			id: "6a30b9b0-18bf-47b4-bd23-d72726864def",
			input: &domain.Task{
				Description: "Update Test Description",
				Status:      true,
			},
			mockSetup: func(m sqlmock.Sqlmock, id string, param *model.UpdateTaskParam) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "status"}).
					AddRow("6a30b9b0-18bf-47b4-bd23-d72726864def", "Test Title", param.Description, param.Status)

				m.ExpectQuery(regexp.QuoteMeta(
					`UPDATE tasks SET description=$1, status=$2 WHERE id=$3
					RETURNING id, title, description, status`,
				)).
					WithArgs(param.Description, param.Status, id).
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
		"NotFound": {
			id: "3e440171-0921-4c88-a7ec-13f4cdab0d69",
			input: &domain.Task{
				Description: "Update Test Description",
				Status:      true,
			},
			mockSetup: func(m sqlmock.Sqlmock, id string, param *model.UpdateTaskParam) {
				sqlmock.NewRows([]string{"id", "title", "description", "status"}).
					AddRow("6a30b9b0-18bf-47b4-bd23-d72726864def", "Test Title", param.Description, param.Status)

				m.ExpectQuery(regexp.QuoteMeta(
					`UPDATE tasks SET description=$1, status=$2 WHERE id=$3
					RETURNING id, title, description, status`,
				)).
					WithArgs(param.Description, param.Status, id).
					WillReturnError(sql.ErrNoRows)
			},
			expected: expected{
				task: nil,
				err:  customError.ErrNotFound,
			},
		},
		"InvalidId": {
			id: "abc123",
			input: &domain.Task{
				Description: "Update Test Description",
				Status:      true,
			},
			mockSetup: func(m sqlmock.Sqlmock, id string, param *model.UpdateTaskParam) {
				sqlmock.NewRows([]string{"id", "title", "description", "status"}).
					AddRow("6a30b9b0-18bf-47b4-bd23-d72726864def", "Test Title", param.Description, param.Status)

				m.ExpectQuery(regexp.QuoteMeta(
					`UPDATE tasks SET description=$1, status=$2 WHERE id=$3
					RETURNING id, title, description, status`,
				)).
					WithArgs(param.Description, param.Status, id).
					WillReturnError(errors.New("pq: invalid input syntax for type uuid: \"abc123\""))
			},
			expected: expected{
				task: nil,
				err:  customError.ErrInternalServerError,
			},
		},
	}

	for n, tt := range testTable {
		tt := tt
		t.Run(n, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock, tt.id, model.ToUpdateTaskParam(tt.input))

			repo := &TaskRepository{Db: db}
			result, err := repo.Update(context.Background(), tt.id, tt.input)

			if tt.expected.err != nil {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expected.err.Error())
			} else {
				assert.Equal(t, tt.expected.task, result)
				assert.Nil(t, err)
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
		"Ok": {
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
		"NotFound": {
			id: "3e440171-0921-4c88-a7ec-13f4cdab0d69",
			mockSetup: func(m sqlmock.Sqlmock, id string) {
				sqlmock.NewRows([]string{"id"}).
					AddRow("6a30b9b0-18bf-47b4-bd23-d72726864def")

				m.ExpectQuery(regexp.QuoteMeta(
					"DELETE FROM tasks WHERE id=$1 RETURNING id",
				)).
					WithArgs(id).
					WillReturnError(sql.ErrNoRows)
			},
			expected: expected{
				task: nil,
				err:  customError.ErrNotFound,
			},
		},
		"InvalidId": {
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
				err:  customError.ErrInternalServerError,
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

			if tt.expected.err != nil {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expected.err.Error())
			} else {
				assert.Equal(t, tt.expected.task, result)
				assert.Nil(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

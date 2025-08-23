package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/takumi616/go-restapi/domain"
	"github.com/takumi616/go-restapi/infrastructure/db/repository/model"
	customError "github.com/takumi616/go-restapi/shared/error"
)

type TaskRepository struct {
	Db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		Db: db,
	}
}

func (r *TaskRepository) Insert(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	param := model.ToInsertTaskParam(task)

	var result model.TaskResult
	err := r.Db.QueryRowContext(
		ctx,
		`INSERT INTO tasks(title, description, status)
		VALUES($1, $2, $3)
		RETURNING id, title, description, status`,
		param.Title, param.Description, param.Status,
	).Scan(&result.Id, &result.Title, &result.Description, &result.Status)

	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, customError.ErrInternalServerError
	}

	return model.ToDomain(&result), nil
}

func (r *TaskRepository) SelectAll(ctx context.Context) ([]*domain.Task, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT id, title, description, status FROM tasks")
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, customError.ErrInternalServerError
	}
	defer rows.Close()

	var taskList []*domain.Task
	for rows.Next() {
		var taskResult model.TaskResult
		if err := rows.Scan(&taskResult.Id, &taskResult.Title, &taskResult.Description, &taskResult.Status); err != nil {
			slog.ErrorContext(ctx, err.Error())
			return nil, customError.ErrInternalServerError
		}
		taskList = append(taskList, model.ToDomain(&taskResult))
	}

	if err := rows.Err(); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, customError.ErrInternalServerError
	}

	if len(taskList) == 0 {
		return []*domain.Task{}, nil
	}

	return taskList, nil
}

func (r *TaskRepository) SelectById(ctx context.Context, id string) (*domain.Task, error) {
	var taskRes model.TaskResult
	err := r.Db.QueryRowContext(
		ctx, "SELECT id, title, description, status FROM tasks WHERE id = $1", id,
	).Scan(&taskRes.Id, &taskRes.Title, &taskRes.Description, &taskRes.Status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, err.Error())
			return nil, customError.ErrNotFound
		}

		slog.ErrorContext(ctx, err.Error())
		return nil, customError.ErrInternalServerError
	}

	return model.ToDomain(&taskRes), nil
}

func (r *TaskRepository) Update(ctx context.Context, id string, task *domain.Task) (*domain.Task, error) {
	param := model.ToUpdateTaskParam(task)

	var result model.TaskResult
	err := r.Db.QueryRowContext(
		ctx,
		`UPDATE tasks SET description=$1, status=$2 WHERE id=$3
		RETURNING id, title, description, status`,
		param.Description, param.Status, id,
	).Scan(&result.Id, &result.Title, &result.Description, &result.Status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, err.Error())
			return nil, customError.ErrNotFound
		}

		slog.ErrorContext(ctx, err.Error())
		return nil, customError.ErrInternalServerError
	}

	return model.ToDomain(&result), nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) (*domain.Task, error) {
	var deletedId string
	err := r.Db.QueryRowContext(
		ctx, "DELETE FROM tasks WHERE id=$1 RETURNING id", id,
	).Scan(&deletedId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, err.Error())
			return nil, customError.ErrNotFound
		}

		slog.ErrorContext(ctx, err.Error())
		return nil, customError.ErrInternalServerError
	}

	task := &domain.Task{}
	task.Id = deletedId

	return task, nil
}

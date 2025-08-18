package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/takumi616/go-restapi/domain"
	"github.com/takumi616/go-restapi/infrastructure/db/repository/model"
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
		"INSERT INTO tasks(title, description, status) VALUES($1, $2, $3) RETURNING *",
		param.Title, param.Description, param.Status,
	).Scan(&result.Id, &result.Title, &result.Description, &result.Status)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return model.ToDomain(&result), nil
}

func (r *TaskRepository) SelectAll(ctx context.Context) ([]*domain.Task, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT * FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var taskList []*domain.Task
	for rows.Next() {
		var taskResult model.TaskResult
		if err := rows.Scan(&taskResult.Id, &taskResult.Title, &taskResult.Description, &taskResult.Status); err != nil {
			return nil, fmt.Errorf("failed to copy columns: %w", err)
		}
		taskList = append(taskList, model.ToDomain(&taskResult))
	}

	return taskList, nil
}

func (r *TaskRepository) SelectById(ctx context.Context, id string) (*domain.Task, error) {
	var taskRes model.TaskResult
	err := r.Db.QueryRowContext(
		ctx, "SELECT * FROM tasks WHERE id = $1", id,
	).Scan(&taskRes.Id, &taskRes.Title, &taskRes.Description, &taskRes.Status)

	if err != nil {
		return nil, fmt.Errorf("failed to copy columns: %w", err)
	}

	return model.ToDomain(&taskRes), nil
}

func (r *TaskRepository) Update(ctx context.Context, id string, task *domain.Task) (*domain.Task, error) {
	param := model.ToUpdateTaskParam(task)

	var result model.TaskResult
	err := r.Db.QueryRowContext(
		ctx,
		"UPDATE tasks SET description=$1, status=$2 WHERE id=$3 RETURNING *",
		param.Description, param.Status, id,
	).Scan(&result.Id, &result.Title, &result.Description, &result.Status)

	if err != nil {
		return nil, fmt.Errorf("failed to copy columns: %w", err)
	}

	return model.ToDomain(&result), nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) (*domain.Task, error) {
	var deletedId string
	err := r.Db.QueryRowContext(
		ctx, "DELETE FROM tasks WHERE id=$1 RETURNING id", id,
	).Scan(&deletedId)

	if err != nil {
		return nil, fmt.Errorf("failed to copy columns: %w", err)
	}

	task := &domain.Task{}
	task.Id = deletedId

	return task, nil
}

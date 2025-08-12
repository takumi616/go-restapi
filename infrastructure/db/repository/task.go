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

	var result model.InsertTaskResult
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

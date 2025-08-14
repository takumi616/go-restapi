package handler

import (
	"context"

	"github.com/takumi616/go-restapi/domain"
)

type TaskUsecase interface {
	AddTask(ctx context.Context, task *domain.Task) (*domain.Task, error)
	GetTaskList(ctx context.Context) ([]*domain.Task, error)
}

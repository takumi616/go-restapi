package usecase

import (
	"context"

	"github.com/takumi616/go-restapi/domain"
)

type TaskGateway interface {
	AddTask(ctx context.Context, task *domain.Task) (*domain.Task, error)
	GetTaskList(ctx context.Context) ([]*domain.Task, error)
	GetTaskById(ctx context.Context, id string) (*domain.Task, error)
	UpdateTask(ctx context.Context, id string, task *domain.Task) (*domain.Task, error)
}

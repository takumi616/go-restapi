package gateway

import (
	"context"

	"github.com/takumi616/go-restapi/domain"
)

type TaskGateway struct {
	repository TaskRepository
}

func NewTaskGateway(repository TaskRepository) *TaskGateway {
	return &TaskGateway{repository: repository}
}

func (g *TaskGateway) AddTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	return g.repository.Insert(ctx, task)
}

func (g *TaskGateway) GetTaskList(ctx context.Context) ([]*domain.Task, error) {
	return g.repository.SelectAll(ctx)
}

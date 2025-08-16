package usecase

import (
	"context"

	"github.com/takumi616/go-restapi/domain"
)

type TaskUsecase struct {
	gateway TaskGateway
}

func NewTaskUsecase(gateway TaskGateway) *TaskUsecase {
	return &TaskUsecase{
		gateway: gateway,
	}
}

func (u *TaskUsecase) AddTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	// Set default status
	task.Status = false

	return u.gateway.AddTask(ctx, task)
}

func (u *TaskUsecase) GetTaskList(ctx context.Context) ([]*domain.Task, error) {
	return u.gateway.GetTaskList(ctx)
}

func (u *TaskUsecase) GetTaskById(ctx context.Context, id string) (*domain.Task, error) {
	return u.gateway.GetTaskById(ctx, id)
}

func (u *TaskUsecase) UpdateTask(ctx context.Context, id string, task *domain.Task) (*domain.Task, error) {
	return u.gateway.UpdateTask(ctx, id, task)
}

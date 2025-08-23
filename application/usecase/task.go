package usecase

import (
	"context"
	"errors"

	"github.com/takumi616/go-restapi/domain"
	customError "github.com/takumi616/go-restapi/shared/error"
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

	task, err := u.gateway.AddTask(ctx, task)
	if err != nil {
		return nil, customError.ErrAddTask
	}

	return task, nil
}

func (u *TaskUsecase) GetTaskList(ctx context.Context) ([]*domain.Task, error) {
	taskList, err := u.gateway.GetTaskList(ctx)
	if err != nil {
		return nil, customError.ErrGetTaskList
	}

	return taskList, nil
}

func (u *TaskUsecase) GetTaskById(ctx context.Context, id string) (*domain.Task, error) {
	task, err := u.gateway.GetTaskById(ctx, id)
	if err != nil {
		if errors.Is(err, customError.ErrNotFound) {
			return nil, customError.ErrTaskNotFound
		} else {
			return nil, customError.ErrGetTaskById
		}
	}

	return task, nil
}

func (u *TaskUsecase) UpdateTask(ctx context.Context, id string, task *domain.Task) (*domain.Task, error) {
	task, err := u.gateway.UpdateTask(ctx, id, task)
	if err != nil {
		if errors.Is(err, customError.ErrNotFound) {
			return nil, customError.ErrTaskNotFound
		} else {
			return nil, customError.ErrUpdateTask
		}
	}

	return task, nil
}

func (u *TaskUsecase) DeleteTask(ctx context.Context, id string) (*domain.Task, error) {
	task, err := u.gateway.DeleteTask(ctx, id)
	if err != nil {
		if errors.Is(err, customError.ErrNotFound) {
			return nil, customError.ErrTaskNotFound
		} else {
			return nil, customError.ErrDeleteTask
		}
	}

	return task, nil
}

package usecase

import (
	"context"

	"github.com/takumi616/go-restapi/domain"
)

type TaskGateway interface {
	AddTask(ctx context.Context, task *domain.Task) (*domain.Task, error)
}

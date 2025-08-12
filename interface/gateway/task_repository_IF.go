package gateway

import (
	"context"

	"github.com/takumi616/go-restapi/domain"
)

type TaskRepository interface {
	Insert(ctx context.Context, task *domain.Task) (*domain.Task, error)
}

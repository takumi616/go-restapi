package gateway

import (
	"context"

	"github.com/takumi616/go-restapi/domain"
)

type TaskRepository interface {
	Insert(ctx context.Context, task *domain.Task) (*domain.Task, error)
	SelectAll(ctx context.Context) ([]*domain.Task, error)
	SelectById(ctx context.Context, id string) (*domain.Task, error)
	Update(ctx context.Context, id string, task *domain.Task) (*domain.Task, error)
	Delete(ctx context.Context, id string) (*domain.Task, error)
}

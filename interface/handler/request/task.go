package request

import "github.com/takumi616/go-restapi/domain"

type AddTaskReq struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (a *AddTaskReq) ToDomain() *domain.Task {
	return &domain.Task{
		Title:       a.Title,
		Description: a.Description,
	}
}

type UpdateTaskReq struct {
	Description string `json:"description" validate:"required"`
	Status      *bool  `json:"status" validate:"required"`
}

func (u *UpdateTaskReq) ToDomain() *domain.Task {
	return &domain.Task{
		Description: u.Description,
		Status:      *u.Status,
	}
}

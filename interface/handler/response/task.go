package response

import "github.com/takumi616/go-restapi/domain"

type TaskRes struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
}

func ToTaskRes(task *domain.Task) *TaskRes {
	return &TaskRes{
		task.Id, task.Title, task.Description, task.Status,
	}
}

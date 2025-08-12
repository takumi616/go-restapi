package response

import "github.com/takumi616/go-restapi/domain"

type AddTaskRes struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
}

func ToAddTaskRes(task *domain.Task) *AddTaskRes {
	return &AddTaskRes{
		task.Id, task.Title, task.Description, task.Status,
	}
}

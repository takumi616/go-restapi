package model

import "github.com/takumi616/go-restapi/domain"

type InsertTaskParam struct {
	Title       string
	Description string
	Status      bool
}

func ToInsertTaskParam(task *domain.Task) *InsertTaskParam {
	return &InsertTaskParam{task.Title, task.Description, task.Status}
}

type TaskResult struct {
	Id          string
	Title       string
	Description string
	Status      bool
}

func ToDomain(result *TaskResult) *domain.Task {
	task := domain.NewTask(result.Title, result.Description)
	task.Id = result.Id
	task.Status = result.Status

	return task
}

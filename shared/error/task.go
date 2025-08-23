package error

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrNotFound            = errors.New("not found")
)

var (
	ErrAddTask      = errors.New("failed to add a new task")
	ErrGetTaskById  = errors.New("failed to get a task by id")
	ErrGetTaskList  = errors.New("failed to get task list")
	ErrUpdateTask   = errors.New("failed to update a task")
	ErrDeleteTask   = errors.New("failed to delete a task")
	ErrTaskNotFound = errors.New("task specified by requested id not found")
)

var (
	TaskBadRequest       = errors.New("requested task info is incorrect")
	InvalidRequestFormat = errors.New("request format is invalid")
)

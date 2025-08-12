package domain

type Task struct {
	Id          string
	Title       string
	Description string
	Status      bool
}

func NewTask(title, description string) *Task {
	return &Task{
		Title:       title,
		Description: description,
	}
}

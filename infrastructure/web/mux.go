package web

import (
	"net/http"

	"github.com/takumi616/go-restapi/interface/handler"
)

type ServeMux struct {
	TaskHandler *handler.TaskHandler
}

func NewServeMux(handler *handler.TaskHandler) *ServeMux {
	return &ServeMux{
		TaskHandler: handler,
	}
}

func (s ServeMux) RegisterHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /tasks", s.TaskHandler.AddTask)
	mux.HandleFunc("GET /tasks", s.TaskHandler.GetTaskList)

	return mux
}

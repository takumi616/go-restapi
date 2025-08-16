package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/takumi616/go-restapi/interface/handler/helper"
	"github.com/takumi616/go-restapi/interface/handler/request"
	"github.com/takumi616/go-restapi/interface/handler/response"
)

type TaskHandler struct {
	usecase TaskUsecase
}

func NewTaskHandler(usecase TaskUsecase) *TaskHandler {
	return &TaskHandler{
		usecase: usecase,
	}
}

func (h *TaskHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req request.AddTaskReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteResponse(
			ctx, w, http.StatusInternalServerError,
			response.ErrResponse{Message: err.Error()},
		)
		return
	}
	defer r.Body.Close()

	err := validator.New().Struct(req)
	if err != nil {
		helper.WriteResponse(
			ctx, w, http.StatusBadRequest,
			response.ErrResponse{Message: err.Error()},
		)
		return
	}

	task := (&req).ToDomain()

	added, err := h.usecase.AddTask(ctx, task)
	if err != nil {
		helper.WriteResponse(
			ctx, w, http.StatusInternalServerError,
			response.ErrResponse{Message: err.Error()},
		)
		return
	}

	helper.WriteResponse(
		ctx, w, http.StatusCreated,
		response.ToTaskRes(added),
	)
}

func (h *TaskHandler) GetTaskList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	taskList, err := h.usecase.GetTaskList(ctx)
	if err != nil {
		helper.WriteResponse(
			ctx, w, http.StatusInternalServerError,
			response.ErrResponse{Message: err.Error()},
		)
		return
	}

	taskResList := []*response.TaskRes{}
	for _, task := range taskList {
		taskResList = append(taskResList, response.ToTaskRes(task))
	}

	helper.WriteResponse(ctx, w, http.StatusOK, taskResList)
}

func (h *TaskHandler) GetTaskById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := r.PathValue("id")
	task, err := h.usecase.GetTaskById(ctx, id)
	if err != nil {
		helper.WriteResponse(
			ctx, w, http.StatusInternalServerError,
			response.ErrResponse{Message: err.Error()},
		)
		return
	}

	helper.WriteResponse(ctx, w, http.StatusOK, response.ToTaskRes(task))
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := r.PathValue("id")
	var req request.UpdateTaskReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.WriteResponse(
			ctx, w, http.StatusInternalServerError,
			response.ErrResponse{Message: err.Error()},
		)
		return
	}
	defer r.Body.Close()

	err := validator.New().Struct(req)
	if err != nil {
		helper.WriteResponse(
			ctx, w, http.StatusBadRequest,
			response.ErrResponse{Message: err.Error()},
		)
		return
	}

	task := (&req).ToDomain()

	updated, err := h.usecase.UpdateTask(ctx, id, task)
	if err != nil {
		helper.WriteResponse(
			ctx, w, http.StatusInternalServerError,
			response.ErrResponse{Message: err.Error()},
		)
		return
	}

	helper.WriteResponse(
		ctx, w, http.StatusOK,
		response.ToTaskRes(updated),
	)
}

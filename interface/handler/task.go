package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/takumi616/go-restapi/interface/handler/helper"
	"github.com/takumi616/go-restapi/interface/handler/request"
	"github.com/takumi616/go-restapi/interface/handler/response"
	customError "github.com/takumi616/go-restapi/shared/error"
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
		slog.ErrorContext(ctx, err.Error())
		helper.WriteResponse(
			ctx, w, http.StatusInternalServerError,
			response.ErrResponse{Message: customError.InvalidRequestFormat.Error()},
		)
		return
	}
	defer r.Body.Close()

	err := validator.New().Struct(req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		helper.WriteResponse(
			ctx, w, http.StatusBadRequest,
			response.ErrResponse{Message: customError.TaskBadRequest.Error()},
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
		if errors.Is(err, customError.ErrTaskNotFound) {
			helper.WriteResponse(
				ctx, w, http.StatusNotFound,
				response.ErrResponse{Message: err.Error()},
			)
		} else {
			helper.WriteResponse(
				ctx, w, http.StatusInternalServerError,
				response.ErrResponse{Message: err.Error()},
			)
		}

		return
	}

	helper.WriteResponse(ctx, w, http.StatusOK, response.ToTaskRes(task))
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := r.PathValue("id")
	var req request.UpdateTaskReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(ctx, err.Error())
		helper.WriteResponse(
			ctx, w, http.StatusInternalServerError,
			response.ErrResponse{Message: customError.InvalidRequestFormat.Error()},
		)
		return
	}
	defer r.Body.Close()

	err := validator.New().Struct(req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		helper.WriteResponse(
			ctx, w, http.StatusBadRequest,
			response.ErrResponse{Message: customError.TaskBadRequest.Error()},
		)
		return
	}

	task := (&req).ToDomain()

	updated, err := h.usecase.UpdateTask(ctx, id, task)
	if err != nil {
		if errors.Is(err, customError.ErrTaskNotFound) {
			helper.WriteResponse(
				ctx, w, http.StatusNotFound,
				response.ErrResponse{Message: err.Error()},
			)
		} else {
			helper.WriteResponse(
				ctx, w, http.StatusInternalServerError,
				response.ErrResponse{Message: err.Error()},
			)
		}

		return
	}

	helper.WriteResponse(
		ctx, w, http.StatusOK,
		response.ToTaskRes(updated),
	)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := r.PathValue("id")
	deleted, err := h.usecase.DeleteTask(ctx, id)
	if err != nil {
		if errors.Is(err, customError.ErrTaskNotFound) {
			helper.WriteResponse(
				ctx, w, http.StatusNotFound,
				response.ErrResponse{Message: err.Error()},
			)
		} else {
			helper.WriteResponse(
				ctx, w, http.StatusInternalServerError,
				response.ErrResponse{Message: err.Error()},
			)
		}

		return
	}

	helper.WriteResponse(ctx, w, http.StatusOK, response.ToTaskIdRes(deleted))
}

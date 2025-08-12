package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/takumi616/go-restapi/domain"
	"github.com/takumi616/go-restapi/interface/handler/helper"
	"github.com/takumi616/go-restapi/interface/handler/request"
	"github.com/takumi616/go-restapi/interface/handler/response"
)

type TaskHandler struct {
	usecase   TaskUsecase
	Validator *validator.Validate
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

	task := domain.NewTask(req.Title, req.Description)

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
		response.ToAddTaskRes(added),
	)
}

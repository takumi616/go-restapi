package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/takumi616/go-restapi/domain"
	"github.com/takumi616/go-restapi/interface/handler/test/helper"
	"github.com/takumi616/go-restapi/interface/handler/test/mock"
	customError "github.com/takumi616/go-restapi/shared/error"
)

func TestAddTask(t *testing.T) {
	type expected struct {
		status  int
		resFile string
	}

	type mockData struct {
		param, returned *domain.Task
		err             error
	}

	testTable := map[string]struct {
		reqFile  string
		expected expected
		mockData mockData
		mockUse  bool
	}{
		"Ok": {
			reqFile: "test/data/add_task/ok_req.json.golden",
			expected: expected{
				status:  http.StatusCreated,
				resFile: "test/data/add_task/ok_res.json.golden",
			},
			mockData: mockData{
				param: &domain.Task{Title: "test title", Description: "test description"},
				returned: &domain.Task{
					Id:    "6a30b9b0-18bf-47b4-bd23-d72726864def",
					Title: "test title", Description: "test description",
					Status: false,
				},
				err: nil,
			},
			mockUse: true,
		},
		"DuplicateErr": {
			reqFile: "test/data/add_task/duplicate_err_req.json.golden",
			expected: expected{
				status:  http.StatusInternalServerError,
				resFile: "test/data/add_task/duplicate_err_res.json.golden",
			},
			mockData: mockData{
				param:    &domain.Task{Title: "duplicate test title", Description: "test description"},
				returned: nil,
				err:      customError.ErrAddTask,
			},
			mockUse: true,
		},
		"UnmarshalFail": {
			reqFile: "test/data/add_task/unmarshal_fail_req.json.golden",
			expected: expected{
				status:  http.StatusInternalServerError,
				resFile: "test/data/add_task/unmarshal_fail_res.json.golden",
			},
			mockData: mockData{
				param:    nil,
				returned: nil,
			},
			mockUse: false,
		},
		"BadRequest": {
			reqFile: "test/data/add_task/bad_req_req.json.golden",
			expected: expected{
				status:  http.StatusBadRequest,
				resFile: "test/data/add_task/bad_req_res.json.golden",
			},
			mockData: mockData{
				param:    nil,
				returned: nil,
			},
			mockUse: false,
		},
	}

	for n, tt := range testTable {
		tt := tt

		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				bytes.NewReader(helper.LoadFile(t, tt.reqFile)),
			)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockTaskUsecase := mock.NewMockTaskUsecase(mockCtrl)
			if tt.mockUse {
				mockTaskUsecase.EXPECT().AddTask(r.Context(), tt.mockData.param).
					Return(tt.mockData.returned, tt.mockData.err)
			}

			sut := NewTaskHandler(mockTaskUsecase)
			sut.AddTask(w, r)

			actualRes := w.Result()
			helper.AssertResponse(t,
				actualRes, tt.expected.status, helper.LoadFile(t, tt.expected.resFile),
			)
		})
	}
}

func TestGetTaskList(t *testing.T) {
	type expected struct {
		status  int
		resFile string
	}

	testTable := map[string]struct {
		taskList []*domain.Task
		err      error
		expected expected
	}{
		"Ok": {
			taskList: []*domain.Task{
				{
					Id:          "f299e7ed-a22a-4494-b59e-21bb91fdae3b",
					Title:       "test title",
					Description: "test description",
					Status:      false,
				},
				{
					Id:          "4d758d63-5c4f-4bef-9a80-d5837c324a07",
					Title:       "test title2",
					Description: "test description2",
					Status:      false,
				},
			},
			err: nil,
			expected: expected{
				status:  http.StatusOK,
				resFile: "test/data/get_task_list/ok_res.json.golden",
			},
		},
		"Empty": {
			taskList: []*domain.Task{},
			err:      nil,
			expected: expected{
				status:  http.StatusOK,
				resFile: "test/data/get_task_list/empty_res.json.golden",
			},
		},
		"InternalServerErr": {
			taskList: nil,
			err:      customError.ErrGetTaskList,
			expected: expected{
				status:  http.StatusInternalServerError,
				resFile: "test/data/get_task_list/internal_server_err_res.json.golden",
			},
		},
	}

	for n, tt := range testTable {
		tt := tt

		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/tasks", nil)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockTaskUsecase := mock.NewMockTaskUsecase(mockCtrl)
			mockTaskUsecase.EXPECT().GetTaskList(r.Context()).
				Return(tt.taskList, tt.err)

			sut := NewTaskHandler(mockTaskUsecase)
			sut.GetTaskList(w, r)

			actualRes := w.Result()
			helper.AssertResponse(t,
				actualRes, tt.expected.status, helper.LoadFile(t, tt.expected.resFile),
			)
		})
	}
}

func TestGetTaskById(t *testing.T) {
	type expected struct {
		status  int
		resFile string
	}

	testTable := map[string]struct {
		id       string
		task     *domain.Task
		err      error
		expected expected
	}{
		"Ok": {
			id: "f299e7ed-a22a-4494-b59e-21bb91fdae3b",
			task: &domain.Task{
				Id:          "f299e7ed-a22a-4494-b59e-21bb91fdae3b",
				Title:       "test title",
				Description: "test description",
				Status:      false,
			},
			err: nil,
			expected: expected{
				status:  http.StatusOK,
				resFile: "test/data/get_task_by_id/ok_res.json.golden",
			},
		},
		"NotFound": {
			id:   "f299e7ed-a22a-4494-b59e-21bb91fdae3b",
			task: nil,
			err:  customError.ErrTaskNotFound,
			expected: expected{
				status:  http.StatusNotFound,
				resFile: "test/data/get_task_by_id/not_found_res.json.golden",
			},
		},
		"InvalidId": {
			id:   "abc123",
			task: nil,
			err:  customError.ErrGetTaskById,
			expected: expected{
				status:  http.StatusInternalServerError,
				resFile: "test/data/get_task_by_id/invalid_id_res.json.golden",
			},
		},
	}

	for n, tt := range testTable {
		tt := tt

		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%s", tt.id), nil)
			r.SetPathValue("id", tt.id)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockTaskUsecase := mock.NewMockTaskUsecase(mockCtrl)
			mockTaskUsecase.EXPECT().GetTaskById(r.Context(), tt.id).
				Return(tt.task, tt.err)

			sut := NewTaskHandler(mockTaskUsecase)
			sut.GetTaskById(w, r)

			actualRes := w.Result()
			helper.AssertResponse(t,
				actualRes, tt.expected.status, helper.LoadFile(t, tt.expected.resFile),
			)
		})
	}
}

func TestUpdateTask(t *testing.T) {
	type expected struct {
		status  int
		resFile string
	}

	type mockData struct {
		inputTask, returnedTask *domain.Task
		err                     error
	}

	testTable := map[string]struct {
		id       string
		reqFile  string
		expected expected
		mockData mockData
		mockUse  bool
	}{
		"Ok": {
			id:      "6a30b9b0-18bf-47b4-bd23-d72726864def",
			reqFile: "test/data/update_task/ok_req.json.golden",
			expected: expected{
				status:  http.StatusOK,
				resFile: "test/data/update_task/ok_res.json.golden",
			},
			mockData: mockData{
				inputTask: &domain.Task{Description: "update test description", Status: true},
				returnedTask: &domain.Task{
					Id:    "6a30b9b0-18bf-47b4-bd23-d72726864def",
					Title: "test title", Description: "update test description",
					Status: true,
				},
				err: nil,
			},
			mockUse: true,
		},
		"UnmarshalFail": {
			reqFile: "test/data/update_task/unmarshal_fail_req.json.golden",
			expected: expected{
				status:  http.StatusInternalServerError,
				resFile: "test/data/update_task/unmarshal_fail_res.json.golden",
			},
			mockData: mockData{
				inputTask:    nil,
				returnedTask: nil,
				err:          nil,
			},
			mockUse: false,
		},
		"BadRequest": {
			reqFile: "test/data/update_task/bad_req_req.json.golden",
			expected: expected{
				status:  http.StatusBadRequest,
				resFile: "test/data/update_task/bad_req_res.json.golden",
			},
			mockData: mockData{
				inputTask:    nil,
				returnedTask: nil,
				err:          nil,
			},
			mockUse: false,
		},
		"NotFound": {
			id:      "7a30b9b0-18bf-47b4-bd23-d72726864def",
			reqFile: "test/data/update_task/not_found_req.json.golden",
			expected: expected{
				status:  http.StatusNotFound,
				resFile: "test/data/update_task/not_found_res.json.golden",
			},
			mockData: mockData{
				inputTask:    &domain.Task{Description: "update test description", Status: true},
				returnedTask: nil,
				err:          customError.ErrTaskNotFound,
			},
			mockUse: true,
		},
		"InvalidId": {
			id:      "7a30b9b0-18bf-47b4-bd23-d72726864def",
			reqFile: "test/data/update_task/invalid_id_req.json.golden",
			expected: expected{
				status:  http.StatusInternalServerError,
				resFile: "test/data/update_task/invalid_id_res.json.golden",
			},
			mockData: mockData{
				inputTask:    &domain.Task{Description: "update test description", Status: true},
				returnedTask: nil,
				err:          customError.ErrUpdateTask,
			},
			mockUse: true,
		},
	}

	for n, tt := range testTable {
		tt := tt

		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPatch,
				fmt.Sprintf("/tasks/%s", tt.id),
				bytes.NewReader(helper.LoadFile(t, tt.reqFile)),
			)
			r.SetPathValue("id", tt.id)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockTaskUsecase := mock.NewMockTaskUsecase(mockCtrl)

			if tt.mockUse {
				mockTaskUsecase.EXPECT().UpdateTask(r.Context(), tt.id, tt.mockData.inputTask).
					Return(tt.mockData.returnedTask, tt.mockData.err)
			}

			sut := NewTaskHandler(mockTaskUsecase)
			sut.UpdateTask(w, r)

			actualRes := w.Result()
			helper.AssertResponse(t,
				actualRes, tt.expected.status, helper.LoadFile(t, tt.expected.resFile),
			)
		})
	}
}

func TestDeleteTask(t *testing.T) {
	type expected struct {
		status  int
		resFile string
	}

	type mockData struct {
		returnedTask *domain.Task
		err          error
	}

	testTable := map[string]struct {
		id       string
		expected expected
		mockData mockData
	}{
		"Ok": {
			id: "f299e7ed-a22a-4494-b59e-21bb91fdae3b",
			expected: expected{
				status:  http.StatusOK,
				resFile: "test/data/delete_task/ok_res.json.golden",
			},
			mockData: mockData{
				returnedTask: &domain.Task{
					Id: "f299e7ed-a22a-4494-b59e-21bb91fdae3b",
				},
				err: nil,
			},
		},
		"NotFound": {
			id: "f299e7ed-a22a-4494-b59e-21bb91fdae3b",
			expected: expected{
				status:  http.StatusNotFound,
				resFile: "test/data/delete_task/not_found_res.json.golden",
			},
			mockData: mockData{
				returnedTask: nil,
				err:          customError.ErrTaskNotFound,
			},
		},
		"InvalidId": {
			id: "abc123",
			expected: expected{
				status:  http.StatusInternalServerError,
				resFile: "test/data/delete_task/invalid_id_res.json.golden",
			},
			mockData: mockData{
				returnedTask: nil,
				err:          customError.ErrDeleteTask,
			},
		},
	}

	for n, tt := range testTable {
		tt := tt

		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%s", tt.id), nil)
			r.SetPathValue("id", tt.id)

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockTaskUsecase := mock.NewMockTaskUsecase(mockCtrl)
			mockTaskUsecase.EXPECT().DeleteTask(r.Context(), tt.id).
				Return(tt.mockData.returnedTask, tt.mockData.err)

			sut := NewTaskHandler(mockTaskUsecase)
			sut.DeleteTask(w, r)

			actualRes := w.Result()
			helper.AssertResponse(t,
				actualRes, tt.expected.status, helper.LoadFile(t, tt.expected.resFile),
			)
		})
	}
}

package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/takumi616/go-restapi/domain"
	"github.com/takumi616/go-restapi/interface/handler/test/helper"
	"github.com/takumi616/go-restapi/interface/handler/test/mock"
)

func TestAddTask(t *testing.T) {
	type expected struct {
		status  int
		resFile string
	}

	type mockData struct {
		param, returned *domain.Task
	}

	testTable := map[string]struct {
		reqFile  string
		expected expected
		mockData mockData
	}{
		"ok": {
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
			},
		},
		"unmarshalFail": {
			reqFile: "test/data/add_task/unmarshal_fail_req.json.golden",
			expected: expected{
				status:  http.StatusInternalServerError,
				resFile: "test/data/add_task/unmarshal_fail_res.json.golden",
			},
			mockData: mockData{
				param:    nil,
				returned: nil,
			},
		},
		"badRequest": {
			reqFile: "test/data/add_task/bad_req.json.golden",
			expected: expected{
				status:  http.StatusBadRequest,
				resFile: "test/data/add_task/bad_res.json.golden",
			},
			mockData: mockData{
				param:    nil,
				returned: nil,
			},
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
			if n == "ok" {
				mockTaskUsecase.EXPECT().AddTask(r.Context(), tt.mockData.param).
					Return(tt.mockData.returned, nil)
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

package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/takumi616/go-restapi/interface/handler/response"
)

func WriteResponse(ctx context.Context, w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Failed to encode response correctly: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		rsp := response.ErrResponse{
			Message: http.StatusText(http.StatusInternalServerError),
		}

		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			fmt.Printf("Failed to write error response correctly: %v", err)
		}
		return
	}

	w.WriteHeader(status)
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		fmt.Printf("Failed to write response correctly: %v", err)
	}
}

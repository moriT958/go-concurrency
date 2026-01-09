package internal

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type AppServer struct{}

func (svr *AppServer) ConvertGrayscalePost(w http.ResponseWriter, r *http.Request) {
	var Url ConvertRequestSchema
	if err := json.NewDecoder(r.Body).Decode(&Url); err != nil {
		msg := "failed to decode request body"
		slog.Error(msg,
			"error", err.Error(),
		)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	if err := Url.validate(); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Enqueue GrayScale Conversion Task
	// ...Asynq enqueue code here...

	resp := ConvertResponseSchema{
		TaskId: "dummyTaskId",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		msg := "failed to encode response body"
		slog.Error(msg,
			"error", err.Error(),
		)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

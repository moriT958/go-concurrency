package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConvertGrayscalePost_正常系(t *testing.T) {
	svr := new(AppServer)

	body := `{"url":"sample-url"}`
	req := httptest.NewRequest(http.MethodPost, "/convert", strings.NewReader(body))
	rec := httptest.NewRecorder()

	// exec handler
	svr.ConvertGrayscalePost(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	// Assert: status 202 Accepted
	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("expected %d, got %d", http.StatusAccepted, resp.StatusCode)
	}

	var got ConvertResponseSchema
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("レスポンスのデコードに失敗: %v", err)
	}

	// Assert: taskId in response body
	want := "dummyTaskId"
	if got.TaskId != want {
		t.Errorf("expected %s, got %s", want, got.TaskId)
	}
}

func TestConvertGrayscalePost_異常系(t *testing.T) {
	testTable := []struct {
		name           string
		requestBody    string
		responseStatus int
	}{
		{"存在しないフィールド名", `{"invalid":"failed"}`, http.StatusBadRequest},
		{"フィールドが空", `{}`, http.StatusBadRequest},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			svr := new(AppServer)

			body := tt.requestBody
			req := httptest.NewRequest(http.MethodPost, "/convert", strings.NewReader(body))
			rec := httptest.NewRecorder()

			// exec handler
			svr.ConvertGrayscalePost(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()

			// Assert: status code
			if resp.StatusCode != tt.responseStatus {
				t.Errorf("expected %d, got %d", tt.responseStatus, resp.StatusCode)
			}
		})
	}
}

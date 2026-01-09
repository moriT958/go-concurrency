package controller

import "errors"

// グレースケール変換のリクエストスキーマ
type ConvertRequestSchema struct {
	Url string `json:"url"`
}

type ConvertResponseSchema struct {
	TaskId string `json:"taskId"`
}

func (req *ConvertRequestSchema) validate() error {
	if req.Url == "" {
		return errors.New("request body url can't be empty")
	}
	return nil
}

// タスクの状態取得
type ConvertTaskResponseSchema struct {
	Status string `json:"status"`
	ObjKey string `json:"objectKey,omitempty"`
}

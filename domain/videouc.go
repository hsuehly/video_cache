package domain

import (
	"context"
)

type UrlBody struct {
	Url string `json:"url"`
}

type DelM4Body struct {
	Url string `json:"url"`
}
type TaskM3Body struct {
	Url    string `json:"url"`
	Surl   string `json:"surl"`
	TaskId int    `json:"taskId"`
	Time   int    `json:"time"`
}
type GetLinkResp struct {
	Url     string `json:"url"`
	Type    string `json:"type"`
	ExpTime int    `json:"expTime"`
}
type TaskM4Body struct {
	Url  string `json:"url"`
	Surl string `json:"surl"`
	Time int    `json:"time"`
}
type DelM3Body struct {
	Url  string `json:"url"`
	Save bool   `json:"save"`
}

type VideoUsecase interface {
	DelM4(c context.Context, url string) error
	DelM3(c context.Context, body *DelM3Body) error
	GetLink(c context.Context, url string) (*GetLinkResp, error)
	TaskM3(c context.Context, body *TaskM3Body) (string, error)
	TaskM4(c context.Context, body *TaskM4Body) error
	SaveM3(c context.Context, url string, data *[]byte) error
}

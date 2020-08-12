package controllers

import (
	"net/http"

	"models/entity"
	"services"
)

func HandlerPostReport(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "topic_post_report")
	if r.Method == http.MethodPost {
		PostPostReport(w, r)
	} else {
		response405(w, r)
	}
}

// PostPostReport
func PostPostReport(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	report := &entity.PostReport{}
	checkRequestBody(w, r, report)
	// 开始插入
	report, err := services.AddPostReport(user.Id, couple.Id, report)
	response417ErrToast(w, r, err)
	// 返回
	response200Data(w, r, struct {
		PostReport *entity.PostReport `json:"postReport"`
	}{report})
}

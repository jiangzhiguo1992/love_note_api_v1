package controllers

import (
	"net/http"

	"models/entity"
	"services"
)

func HandlerPostCommentReport(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "topic_post_comment_report")
	if r.Method == http.MethodPost {
		PostPostCommentReport(w, r)
	} else {
		response405(w, r)
	}
}

// PostPostCommentReport
func PostPostCommentReport(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	report := &entity.PostCommentReport{}
	checkRequestBody(w, r, report)
	// 开始插入
	report, err := services.AddPostCommentReport(user.Id, couple.Id, report)
	response417ErrToast(w, r, err)
	// 返回
	response200Data(w, r, struct {
		PostCommentReport *entity.PostCommentReport `json:"postCommentReport"`
	}{report})
}

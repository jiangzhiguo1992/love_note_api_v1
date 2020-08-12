package controllers

import (
	"net/http"

	"models/entity"
	"services"
)

func HandlerPostCommentPoint(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "topic_post_comment_point")
	if r.Method == http.MethodPost {
		PostPostCommentPoint(w, r)
	} else {
		response405(w, r)
	}
}

// PostPostCommentPoint
func PostPostCommentPoint(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	point := &entity.PostCommentPoint{}
	checkRequestBody(w, r, point)
	// 开始插入
	point, err := services.TogglePostCommentPoint(user.Id, couple.Id, point)
	response417ErrToast(w, r, err)
	// 返回
	response200Data(w, r, struct {
		PostCommentPoint *entity.PostCommentPoint `json:"postCommentPoint"`
	}{point})
}

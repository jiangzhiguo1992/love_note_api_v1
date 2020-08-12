package controllers

import (
	"net/http"

	"models/entity"
	"services"
)

func HandlerPostPoint(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "topic_post_point")
	if r.Method == http.MethodPost {
		PostPostPoint(w, r)
	} else {
		response405(w, r)
	}
}

// PostPostPoint
func PostPostPoint(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	point := &entity.PostPoint{}
	checkRequestBody(w, r, point)
	// 开始插入
	point, err := services.TogglePostPoint(user.Id, couple.Id, point)
	response417ErrToast(w, r, err)
	// 返回
	response200Data(w, r, struct {
		PostPoint *entity.PostPoint `json:"postPoint"`
	}{point})
}

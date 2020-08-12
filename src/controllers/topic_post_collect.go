package controllers

import (
	"net/http"

	"models/entity"
	"services"
)

func HandlerPostCollect(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "topic_post_collect")
	if r.Method == http.MethodPost {
		PostPostCollect(w, r)
	} else {
		response405(w, r)
	}
}

// PostPostCollect
func PostPostCollect(w http.ResponseWriter, r *http.Request) {
	user, _ := getTokenCouple(r)
	couple := user.Couple
	// 接受参数
	collect := &entity.PostCollect{}
	checkRequestBody(w, r, collect)
	// 开始插入
	collect, err := services.TogglePostCollect(user.Id, couple.Id, collect)
	response417ErrToast(w, r, err)
	// 返回
	response200Data(w, r, struct {
		PostCollect *entity.PostCollect `json:"postCollect"`
	}{collect})
}

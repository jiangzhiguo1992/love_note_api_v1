package controllers

import (
	"models/entity"
	"net/http"
	"services"
)

func HandlerSuggestFollow(w http.ResponseWriter, r *http.Request) {
	checkModelStatus(w, r, "suggest_follow")
	if r.Method == http.MethodPost {
		PostSuggestFollow(w, r)
	} else {
		response405(w, r)
	}
}

// PostSuggestFollow
func PostSuggestFollow(w http.ResponseWriter, r *http.Request) {
	user := getTokenUser(r)
	// 接受参数
	follow := &entity.SuggestFollow{}
	checkRequestBody(w, r, follow)
	follow.UserId = user.Id
	// 开始点赞
	follow, err := services.ToggleSuggestFollow(follow)
	response417ErrToast(w, r, err)
	// 返回
	response200Data(w, r, struct {
		SuggestFollow *entity.SuggestFollow `json:"suggestFollow"`
	}{follow})
}
